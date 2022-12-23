package namenode

import (
	"aDrive/pkg/tree"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/tidwall/wal"

	namenode_pb "aDrive/proto/namenode"
	"time"
)

type NameNodeMetaData struct {
	BlockId        string
	BlockAddresses []DataNodeInstance
}

type ReDistributeDataRequest struct {
	DataNodeUri string
}

type UnderReplicatedBlocks struct {
	BlockId           string
	HealthyDataNodeId int64
}

type DataNodeInstance struct {
	Host        string
	ServicePort string
}

type Service struct {
	namenode_pb.UnimplementedNameNodeServiceServer

	Port                uint16
	IdToDataNodes       map[int64]DataNodeInstance
	FileNameToDataNodes map[string][]DataMessage
	DataNodeMessageMap  map[string]DataNodeMessage
	DataNodeHeartBeat   map[string]time.Time
	DirTree             *tree.DirTree
	RaftNode            *raft.Raft
	RaftLog             *boltdb.BoltStore
	Log                 *wal.Log
}

type DataMessage struct {
	Host string
	Id   int64
}

type DataNodeMessage struct {
	UsedDisk   uint64
	UsedMem    uint64
	CpuPercent float32
	TotalMem   uint64
	TotalDisk  uint64
}

func NewService(r *raft.Raft, log *boltdb.BoltStore, serverPort uint16) *Service {
	return &Service{
		RaftNode:            r,
		RaftLog:             log,
		Port:                serverPort,
		IdToDataNodes:       make(map[int64]DataNodeInstance),
		FileNameToDataNodes: make(map[string][]DataMessage),
		DataNodeMessageMap:  make(map[string]DataNodeMessage),
		DataNodeHeartBeat:   make(map[string]time.Time),
		DirTree:             initDirTree(),
	}
}

func initDirTree() *tree.DirTree {
	root := &tree.DirTreeNode{
		Name:     "/",
		Children: []*tree.DirTreeNode{},
	}
	return &tree.DirTree{Root: root}
}
