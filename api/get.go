package api

import (
	nn "aDrive/proto/namenode"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GetResp struct {
	Data []byte
}

func Get(addr string, absolutePath string) (GetResp, error) {
	clientConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	getResp, err := nn.NewNameNodeServiceClient(clientConn).Get(context.Background(), &nn.GetReq{
		AbsolutePath: absolutePath,
	})
	if err != nil {
		return GetResp{}, err
	}
	return GetResp{
		Data: getResp.Data,
	}, err
}
