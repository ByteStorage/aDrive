package api

import (
	nn "aDrive/proto/namenode"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MkdirResp struct {
}

func Mkdir(addr string, absolutePath string) (MkdirResp, error) {
	clientConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	_, err = nn.NewNameNodeServiceClient(clientConn).Mkdir(context.Background(), &nn.MkdirReq{
		Path: absolutePath,
	})
	if err != nil {
		return MkdirResp{}, err
	}
	return MkdirResp{}, err
}
