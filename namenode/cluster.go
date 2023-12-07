package namenode

import (
	nn "aDrive/proto/namenode"
	"context"
	"encoding/json"
	"errors"
	"github.com/hashicorp/raft"
	"log"
	"time"
)

func (s *Service) FindLeader(c context.Context, req *nn.FindLeaderReq) (*nn.FindLeaderResp, error) {
	id, _ := s.RaftNode.LeaderWithID()
	if id == "" {
		return &nn.FindLeaderResp{}, errors.New("cannot find leader")
	}
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println("cannot marshal data")
		return &nn.FindLeaderResp{}, err
	}
	if s.RaftNode != nil {
		s.RaftNode.Apply(bytes, time.Second*1)
	}
	return &nn.FindLeaderResp{
		Addr: string(id),
	}, nil
}

func (s *Service) JoinCluster(c context.Context, req *nn.JoinClusterReq) (*nn.JoinClusterResp, error) {
	log.Println("申请加入集群的节点信息为:", req.Id, " ", req.Addr)
	voter := s.RaftNode.AddVoter(raft.ServerID(req.Id), raft.ServerAddress(req.Addr), req.PreviousIndex, 0)
	if voter.Error() != nil {
		return &nn.JoinClusterResp{}, voter.Error()
	}
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println("cannot marshal data")
		return &nn.JoinClusterResp{}, err
	}
	if s.RaftNode != nil {
		s.RaftNode.Apply(bytes, time.Second*1)
	}
	return &nn.JoinClusterResp{Success: true}, nil
}
