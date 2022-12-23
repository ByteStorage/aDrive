package namenode

import (
	"aDrive/pkg/utils"
	dn "aDrive/proto/datanode"
	nn "aDrive/proto/namenode"
	"context"
	"errors"
	"github.com/klauspost/reedsolomon"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
)

func (s *Service) Put(c context.Context, req *nn.PutReq) (*nn.PutResp, error) {
	//判断当前map有多少个DataNodes
	num := len(s.IdToDataNodes)
	//采用EC编码将数据切块
	encoder, err := reedsolomon.New(num-num/3, num/3)
	if err != nil {
		zap.L().Error("create ec encoder error" + err.Error())
		return &nn.PutResp{}, errors.New("create ec encoder error")
	}
	split, err := encoder.Split(req.Data)
	if err != nil {
		zap.L().Error("split data error" + err.Error())
		return &nn.PutResp{}, errors.New("split data error")
	}
	err = encoder.Encode(split)
	if err != nil {
		zap.L().Error("encode data error" + err.Error())
		return &nn.PutResp{}, errors.New("encode data error")
	}
	dataMessage := make([]*nn.DataMessage, num)
	//将切块后的数据分别发送给DataNode
	for i := 0; i < num; i++ {
		//将数据发送给DataNode
		dataNodeInstance := s.IdToDataNodes[int64(i)]
		dnConnect, err := grpc.Dial(dataNodeInstance.Host+":"+dataNodeInstance.ServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Error("connect to datanode(" + dataNodeInstance.Host + ":" + dataNodeInstance.ServicePort + "):" + err.Error())
			return &nn.PutResp{}, errors.New("connect to datanode(" + dataNodeInstance.Host + ":" + dataNodeInstance.ServicePort + "):" + err.Error())
		}
		dnClient := dn.NewDataNodeClient(dnConnect)
		putResp, err := dnClient.Put(context.Background(), &dn.PutReq{
			Data:         split[i],
			AbsolutePath: req.AbsolutePath + strconv.Itoa(i),
		})
		if err != nil {
			zap.L().Error("put data to datanode(" + dataNodeInstance.Host + ":" + dataNodeInstance.ServicePort + "):" + err.Error())
			return &nn.PutResp{}, errors.New("put data to datanode(" + dataNodeInstance.Host + ":" + dataNodeInstance.ServicePort + "):" + err.Error())
		}
		if !putResp.Success {
			zap.L().Error("put data to datanode(" + dataNodeInstance.Host + ":" + dataNodeInstance.ServicePort + "):" + err.Error())
			return &nn.PutResp{}, errors.New("put data to datanode(" + dataNodeInstance.Host + ":" + dataNodeInstance.ServicePort + "):" + err.Error())
		}
		s.FileNameToDataNodes[req.AbsolutePath] = append(s.FileNameToDataNodes[req.AbsolutePath], DataMessage{
			Id:   int64(i),
			Host: dataNodeInstance.Host + ":" + dataNodeInstance.ServicePort,
		})
		dataMessage[i] = &nn.DataMessage{
			Id:   int64(i),
			Host: dataNodeInstance.Host + ":" + dataNodeInstance.ServicePort,
		}
	}
	//更新元数据信息
	path := utils.ModPath(req.AbsolutePath)
	ok := s.DirTree.Insert(path)
	if !ok {
		zap.L().Error("insert dir tree path error")
		return &nn.PutResp{}, errors.New("insert dir tree path error")
	}
	return &nn.PutResp{
		AbsolutePath: req.AbsolutePath,
		DataMessage:  dataMessage,
	}, nil
}
