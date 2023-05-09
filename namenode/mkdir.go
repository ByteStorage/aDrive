package namenode

import (
	"aDrive/pkg/utils"
	dn "aDrive/proto/datanode"
	nn "aDrive/proto/namenode"
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func (s *Service) Mkdir(c context.Context, req *nn.MkdirReq) (*nn.MkdirResp, error) {
	path := utils.ModPath(req.Path)
	ok := s.DirTree.Insert(path)
	if !ok {
		zap.S().Error("插入目录失败，请确认操作是否有误")
		return &nn.MkdirResp{}, errors.New("插入目录失败，请确认操作是否有误")
	}
	ok = s.DirTree.Insert(path + "../")
	if !ok {
		zap.S().Error("插入目录失败，请确认操作是否有误")
		return &nn.MkdirResp{}, errors.New("插入目录失败，请确认操作是否有误")
	}

	for i := range s.IdToDataNodes {
		dataNodeInstance := s.IdToDataNodes[i]
		conn, err := grpc.Dial(dataNodeInstance.Host+":"+dataNodeInstance.ServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.S().Error("connect to datanode(" + s.IdToDataNodes[i].Host + "):" + err.Error())
			return &nn.MkdirResp{}, err
		}
		client := dn.NewDataNodeClient(conn)
		_, err = client.Mkdir(context.Background(), &dn.MkdirReq{
			Path: path,
		})
		if err != nil {
			zap.S().Error("mkdir(" + path + ") error:" + err.Error())
			return &nn.MkdirResp{}, err
		}
	}

	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println("cannot marshal data")
		return &nn.MkdirResp{}, err
	}
	s.RaftNode.Apply(bytes, time.Second*1)
	return &nn.MkdirResp{}, nil
}
