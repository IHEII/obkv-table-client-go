package route

import (
	"github.com/oceanbase/obkv-table-client-go/protocol"
	"strconv"
)

type ObReplicaLocation struct {
	addr        ObServerAddr
	info        ObServerInfo
	role        ObServerRole
	replicaType ObReplicaType
}

func (l *ObReplicaLocation) isValid() bool {
	return !l.role.isInvalid() && l.info.isActive()
}

func (l *ObReplicaLocation) isLeader() bool {
	return l.role.isLeader()
}

func (l *ObReplicaLocation) ToString() string {
	return "ObReplicaLocation{" +
		"addr:" + l.addr.ToString() + ", " +
		"info:" + l.info.ToString() + ", " +
		"role:" + l.role.ToString() + ", " +
		"replicaType:" + l.replicaType.ToString() +
		"}"
}

type ObTableLocation struct {
	replicaLocations []ObReplicaLocation
}

func (l *ObTableLocation) ToString() string {
	var replicaLocationsStr string
	replicaLocationsStr = replicaLocationsStr + "["
	for i := 0; i < len(l.replicaLocations); i++ {
		if i > 0 {
			replicaLocationsStr += ", "
		}
		replicaLocationsStr += l.replicaLocations[i].ToString()
	}
	replicaLocationsStr += "]"
	return "ObTableLocation{" +
		"replicaLocations:" + replicaLocationsStr +
		"}"
}

type ObPartitionInfo struct {
	level           ObPartitionLevel
	firstPartDesc   ObPartDesc
	subPartDesc     ObPartDesc
	partColumns     []protocol.ObColumn
	partTabletIdMap map[int64]int64
	partNameIdMap   map[string]int64
	rowKeyElement   map[string]int
}

func (p *ObPartitionInfo) ToString() string {
	// partColumns to string
	var partColumnsStr string
	partColumnsStr = partColumnsStr + "["
	for i := 0; i < len(p.partColumns); i++ {
		if i > 0 {
			partColumnsStr += ", "
		}
		partColumnsStr += p.partColumns[i].ToString()
	}
	partColumnsStr += "]"

	// partTabletIdMap to string
	var partTabletIdMapStr string
	var i = 0
	partTabletIdMapStr = partTabletIdMapStr + "{"
	for k, v := range p.partTabletIdMap {
		if i > 0 {
			partTabletIdMapStr += ", "
		}
		i++
		partTabletIdMapStr += "m[" + strconv.Itoa(int(k)) + "]=" + strconv.Itoa(int(v))
	}
	partTabletIdMapStr += "}"

	// partNameIdMap to string
	var partNameIdMapStr string
	i = 0
	partNameIdMapStr = partNameIdMapStr + "{"
	for k, v := range p.partNameIdMap {
		if i > 0 {
			partNameIdMapStr += ", "
		}
		i++
		partNameIdMapStr += "m[" + k + "]=" + strconv.Itoa(int(v))
	}
	partNameIdMapStr += "}"

	// rowKeyElement to string
	var rowKeyElementStr string
	rowKeyElementStr = rowKeyElementStr + "{"
	i = 0
	for k, v := range p.rowKeyElement {
		if i > 0 {
			rowKeyElementStr += ", "
		}
		i++
		rowKeyElementStr += "m[" + k + "]=" + strconv.Itoa(v)
	}
	rowKeyElementStr += "}"

	// firstPartDesc to string
	var firstPartDescStr string
	if p.firstPartDesc != nil {
		firstPartDescStr = p.firstPartDesc.ToString()
	} else {
		firstPartDescStr = "nil"
	}

	// subPartDesc to string
	var subPartDescStr string
	if p.subPartDesc != nil {
		subPartDescStr = p.firstPartDesc.ToString()
	} else {
		subPartDescStr = "nil"
	}

	return "ObPartitionInfo{" +
		"level:" + p.level.ToString() + ", " +
		"firstPartDesc:" + firstPartDescStr + ", " +
		"subPartDesc:" + subPartDescStr + ", " +
		"partColumns:" + partColumnsStr + ", " +
		"partTabletIdMap:" + partTabletIdMapStr + ", " +
		"partNameIdMap:" + partNameIdMapStr + ", " +
		"rowKeyElement:" + rowKeyElementStr +
		"}"
}

type ObPartitionLocation struct {
	leader   ObReplicaLocation
	replicas []ObReplicaLocation
}

func (l *ObPartitionLocation) addReplicaLocation(replica *ObReplicaLocation) {
	if replica.isLeader() {
		l.leader = *replica
	}
	l.replicas = append(l.replicas, *replica)
}

func (l *ObPartitionLocation) ToString() string {
	var replicasStr string
	replicasStr = replicasStr + "["
	for i := 0; i < len(l.replicas); i++ {
		if i > 0 {
			replicasStr += ", "
		}
		replicasStr += l.replicas[i].ToString()
	}
	replicasStr += "]"
	return "ObPartitionLocation{" +
		"leader:" + l.leader.ToString() + ", " +
		"replicas:" + replicasStr +
		"}"
}

type ObPartLocationEntry struct {
	partLocations map[int]*ObPartitionLocation
}

func newObPartLocationEntry(partNum int) *ObPartLocationEntry {
	entry := new(ObPartLocationEntry)
	entry.partLocations = make(map[int]*ObPartitionLocation, partNum)
	return entry
}

func (e *ObPartLocationEntry) ToString() string {
	var partitionLocationStr string
	var i = 0
	partitionLocationStr = partitionLocationStr + "{"
	for k, v := range e.partLocations {
		if i > 0 {
			partitionLocationStr += ", "
		}
		i++
		partitionLocationStr += "m[" + strconv.Itoa(k) + "]=" + v.ToString()
	}
	partitionLocationStr += "}"
	return "ObPartLocationEntry{" +
		"partLocations:" + partitionLocationStr +
		"}"
}

type ObTableEntryKey struct {
	clusterName  string
	tenantName   string
	databaseName string
	tableName    string
}

func (k *ObTableEntryKey) ToString() string {
	return "ObTableEntryKey{" +
		"clusterName:" + k.clusterName + ", " +
		"tenantNane:" + k.databaseName + ", " +
		"databaseName:" + k.databaseName + ", " +
		"tableName:" + k.tableName +
		"}"
}

type ObTableEntry struct {
	tableId           uint64
	partNum           int
	replicaNum        int
	refreshTimeMills  int64
	tableEntryKey     ObTableEntryKey
	partitionInfo     *ObPartitionInfo
	tableLocation     *ObTableLocation
	partLocationEntry *ObPartLocationEntry
}

func (e *ObTableEntry) IsPartitionTable() bool {
	return e.partNum > 1
}

func (e *ObTableEntry) ToString() string {
	// partitionInfo to string
	var partitionInfoStr string
	if e.partitionInfo != nil {
		partitionInfoStr = e.partitionInfo.ToString()
	} else {
		partitionInfoStr = "nil"
	}

	// tableLocation to string
	var tableLocationStr string
	if e.tableLocation != nil {
		tableLocationStr = e.tableLocation.ToString()
	} else {
		tableLocationStr = "nil"
	}

	// partLocationEntry to string
	var partLocationEntryStr string
	if e.partLocationEntry != nil {
		partLocationEntryStr = e.partLocationEntry.ToString()
	} else {
		partLocationEntryStr = "nil"
	}
	return "ObTableEntry{" +
		"tableId:" + strconv.Itoa(int(e.tableId)) + ", " +
		"partNum:" + strconv.Itoa(int(e.partNum)) + ", " +
		"replicaNum:" + strconv.Itoa(int(e.replicaNum)) + ", " +
		"refreshTimeMills:" + strconv.Itoa(int(e.refreshTimeMills)) + ", " +
		"tableEntryKey:" + e.tableEntryKey.ToString() + ", " +
		"partitionInfo:" + partitionInfoStr + ", " +
		"tableLocation:" + tableLocationStr + ", " +
		"partitionEntry:" + partLocationEntryStr +
		"}"
}