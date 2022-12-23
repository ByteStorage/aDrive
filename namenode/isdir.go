package namenode

import (
	"aDrive/pkg/utils"
	namenode_pb "aDrive/proto/namenode"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"time"
)

func (s *Service) IsDir(c context.Context, req *namenode_pb.IsDirReq) (*namenode_pb.IsDirResp, error) {
	filename := utils.ModPath(req.Filename)
	//空文件夹下面会有..文件夹
	dir := s.DirTree.FindSubDir(filename)
	bytes, err := json.Marshal(s)
	if err != nil {
		zap.L().Error("cannot marshal data:" + err.Error())
		return &namenode_pb.IsDirResp{}, err
	}
	s.RaftNode.Apply(bytes, time.Second*1)
	if len(dir) == 0 {
		return &namenode_pb.IsDirResp{Ok: false}, nil
	} else {
		return &namenode_pb.IsDirResp{Ok: true}, nil
	}
}
