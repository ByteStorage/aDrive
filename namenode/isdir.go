package namenode

import (
	"aDrive/pkg/utils"
	nn "aDrive/proto/namenode"
	"context"
	"encoding/json"
	"log"
	"time"
)

func (s *Service) IsDir(c context.Context, req *nn.IsDirReq) (*nn.IsDirResp, error) {
	filename := utils.ModPath(req.Filename)
	//空文件夹下面会有..文件夹
	dir := s.DirTree.FindSubDir(filename)
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println("cannot marshal data")
		return &nn.IsDirResp{}, err
	}
	if s.RaftNode != nil {
		s.RaftNode.Apply(bytes, time.Second*1)
	}
	if len(dir) == 0 {
		return &nn.IsDirResp{Ok: false}, nil
	} else {
		return &nn.IsDirResp{Ok: true}, nil
	}
}
