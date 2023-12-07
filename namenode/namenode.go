package namenode

import (
	"aDrive/pkg/tree"
	"context"
	"encoding/json"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/tidwall/wal"
	"go.uber.org/zap"
	"log"
	"net"

	nn "aDrive/proto/namenode"
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
	nn.UnimplementedNameNodeServiceServer

	Port                uint16
	IdToDataNodes       map[int64]DataNodeInstance
	FileNameToDataNodes map[string][]DataMessage
	DataNodeMessageMap  map[string]DataNodeMessage
	DataNodeHeartBeat   map[string]time.Time
	IdToData            map[int]int64
	DirTree             *tree.DirTree
	RaftNode            *raft.Raft
	RaftLog             *boltdb.BoltStore
	Log                 *wal.Log
	Peer                []string
	Leader              string
}

type DataMessage struct {
	Host string
	Id   int
}

type DataNodeMessage struct {
	UsedDisk   uint64
	UsedMem    uint64
	CpuPercent float32
	TotalMem   uint64
	TotalDisk  uint64
	Place      string
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
		IdToData:            make(map[int]int64),
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

func (s *Service) HeartBeat(c context.Context, req *nn.HeartBeatReq) (*nn.HeartBeatResp, error) {
	s.DataNodeHeartBeat[req.Addr] = time.Now()
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println("cannot marshal data")
		return &nn.HeartBeatResp{}, err
	}
	if s.RaftNode != nil {
		s.RaftNode.Apply(bytes, time.Second*1)
	}
	return &nn.HeartBeatResp{Success: true}, nil
}

func (s *Service) RegisterDataNode(c context.Context, req *nn.RegisterDataNodeReq) (*nn.RegisterDataNodeResp, error) {
	s.DataNodeHeartBeat[req.Addr] = time.Now()
	s.DataNodeMessageMap[req.Addr] = DataNodeMessage{
		UsedDisk:   req.UsedDisk,
		UsedMem:    req.UsedMem,
		TotalMem:   req.TotalMem,
		TotalDisk:  req.TotalDisk,
		CpuPercent: req.CpuPercent,
		Place:      req.Place,
	}
	host, port, err := net.SplitHostPort(req.Addr)
	s.IdToDataNodes[time.Now().Unix()] = DataNodeInstance{
		Host:        host,
		ServicePort: port,
	}
	if err != nil {
		zap.L().Error("cannot split host port", zap.Error(err))
		return &nn.RegisterDataNodeResp{}, err
	}
	return &nn.RegisterDataNodeResp{Success: true}, nil
}

func (s *Service) UpdateDataNodeMessage(c context.Context, req *nn.UpdateDataNodeMessageReq) (*nn.UpdateDataNodeMessageResp, error) {
	s.DataNodeMessageMap[req.Addr] = DataNodeMessage{
		UsedDisk:   req.UsedDisk,
		UsedMem:    req.UsedMem,
		TotalMem:   req.TotalMem,
		TotalDisk:  req.TotalDisk,
		CpuPercent: req.CpuPercent,
		Place:      req.Place,
	}
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println("cannot marshal data")
		return &nn.UpdateDataNodeMessageResp{}, err
	}
	if s.RaftNode != nil {
		s.RaftNode.Apply(bytes, time.Second*1)
	}
	return &nn.UpdateDataNodeMessageResp{Success: true}, nil
}
