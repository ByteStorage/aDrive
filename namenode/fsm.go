package namenode

import (
	"aDrive/pkg/tree"
	"encoding/json"
	"github.com/hashicorp/raft"
	"io"
	"io/ioutil"
	"path/filepath"
	"time"
)

var _ raft.FSM = &Service{}

func (nameNode *Service) Apply(l *raft.Log) interface{} {

	filePath := filepath.Join("./Cluster", "metadata.dat")

	err := ioutil.WriteFile(filePath, l.Data, 0755)

	//file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	//_, err = file.Write(l.Data)

	if err != nil {
		panic(err)
	}

	return nil

}

// Snapshot 生成快照
func (nameNode *Service) Snapshot() (raft.FSMSnapshot, error) {
	//TODO implement me
	return Snap{
		Port:                nameNode.Port,
		IdToDataNodes:       nameNode.IdToDataNodes,
		DirTree:             nameNode.DirTree,
		DataNodeMessageMap:  nameNode.DataNodeMessageMap,
		FileNameToDataNodes: nameNode.FileNameToDataNodes,
		DataNodeHeartBeat:   nameNode.DataNodeHeartBeat,
	}, nil
}

// Restore 从快照中恢复数据
func (nameNode *Service) Restore(snapshot io.ReadCloser) error {
	//TODO implement me
	var s Snap
	err := json.NewDecoder(snapshot).Decode(&s)
	if err != nil {
		panic(err)
	}
	nameNode.IdToDataNodes = s.IdToDataNodes
	nameNode.FileNameToDataNodes = s.FileNameToDataNodes
	nameNode.DataNodeHeartBeat = s.DataNodeHeartBeat
	nameNode.DirTree = s.DirTree
	nameNode.DataNodeMessageMap = s.DataNodeMessageMap
	return nil
}

// Snap 用于生成快照，服务重启的时候会从快照中恢复，需要实现两个接口
type Snap struct {
	Port                uint16
	BlockSize           uint64
	ReplicationFactor   uint64
	IdToDataNodes       map[int64]DataNodeInstance
	FileNameToDataNodes map[string][]DataMessage
	BlockToDataNodeIds  map[string][]int64
	DataNodeMessageMap  map[string]DataNodeMessage
	DataNodeHeartBeat   map[string]time.Time
	DirTree             *tree.DirTree
}

func (s Snap) Persist(sink raft.SnapshotSink) error {
	//TODO implement me
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = sink.Write(bytes)
	if err != nil {
		return err
	}
	err = sink.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s Snap) Release() {
}
