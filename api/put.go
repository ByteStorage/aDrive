package api

import (
	nn "aDrive/proto/namenode"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PutResp struct {
	AbsolutePath string
	DataMessage  []*nn.DataMessage
}

// Put 给定要发起grpc调用的NameNode地址以及传递参数，来作为api接口
func Put(addr string, absolutePath string, data []byte) (PutResp, error) {
	clientConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return PutResp{}, errors.New("grpc dial" + addr + " error" + err.Error())
	}
	client := nn.NewNameNodeServiceClient(clientConn)
	putResp, err := client.Put(context.Background(), &nn.PutReq{
		AbsolutePath: absolutePath,
		Data:         data,
	})
	if err != nil {
		return PutResp{}, errors.New("grpc put error" + err.Error())
	}
	return PutResp{
		AbsolutePath: putResp.AbsolutePath,
		DataMessage:  putResp.DataMessage,
	}, nil
}
