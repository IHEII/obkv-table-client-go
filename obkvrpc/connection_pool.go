package obkvrpc

import (
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type PoolOption struct {
	ip                  string
	port                int
	connPoolMaxConnSize int
	connectTimeout      time.Duration

	tenantName   string
	databaseName string
	userName     string
	password     string
}

type ConnectionPool struct {
	option *PoolOption

	connections []*Connection
	rwMutexes   []sync.RWMutex
}

func NewPoolOption(ip string, port int, connPoolMaxConnSize int, connectTimeout time.Duration, tenantName string, databaseName string, userName string, password string) *PoolOption {
	return &PoolOption{ip: ip, port: port, connPoolMaxConnSize: connPoolMaxConnSize, connectTimeout: connectTimeout, tenantName: tenantName, databaseName: databaseName, userName: userName, password: password}
}

func NewConnectionPool(option *PoolOption) (*ConnectionPool, error) {
	pool := &ConnectionPool{
		option:      option,
		connections: make([]*Connection, 0, option.connPoolMaxConnSize),
		rwMutexes:   make([]sync.RWMutex, 0, option.connPoolMaxConnSize),
	}

	connectionOption := NewConnectionOption(pool.option.ip, pool.option.port, pool.option.connectTimeout, pool.option.tenantName, pool.option.databaseName, pool.option.userName, pool.option.password)

	for i := 0; i < pool.option.connPoolMaxConnSize; i++ {

		id := uuid.New()
		connection := NewConnection(connectionOption, id)
		err := connection.Connect()
		if err != nil {
			return nil, errors.Wrap(err, "create connection pool failed")
		}

		err = connection.Login()
		if err != nil {
			return nil, errors.Wrap(err, "create connection pool failed")
		}

		pool.connections = append(pool.connections, connection)
		pool.rwMutexes = append(pool.rwMutexes, sync.RWMutex{})

	}

	return pool, nil
}

func (p *ConnectionPool) GetConnection() (*Connection, error) {
	randomIndex := rand.Intn(len(p.connections))

	p.rwMutexes[randomIndex].RLock()
	if p.connections[randomIndex].active.Load() {
		p.rwMutexes[randomIndex].RUnlock()
		return p.connections[randomIndex], nil
	}
	p.rwMutexes[randomIndex].RUnlock()

	p.rwMutexes[randomIndex].Lock()
	if p.connections[randomIndex].active.Load() {
		p.rwMutexes[randomIndex].Unlock()
		return p.connections[randomIndex], nil
	}
	// Recreate the connection and login
	connection, err := p.CreateConnection()
	if err != nil {
		p.rwMutexes[randomIndex].Unlock()
		return nil, errors.Wrap(err, "get connection recreate failed")
	}

	p.connections[randomIndex] = connection

	p.rwMutexes[randomIndex].Unlock()
	return p.connections[randomIndex], nil
}

func (p *ConnectionPool) CreateConnection() (*Connection, error) {
	connectionOption := NewConnectionOption(p.option.ip, p.option.port, p.option.connectTimeout, p.option.tenantName, p.option.databaseName, p.option.userName, p.option.password)
	connection := NewConnection(connectionOption, uuid.New())
	err := connection.Connect()
	if err != nil {
		return nil, errors.Wrap(err, "create connection failed")
	}
	err = connection.Login()
	if err != nil {
		return nil, errors.Wrap(err, "create connection failed")
	}
	return connection, nil
}