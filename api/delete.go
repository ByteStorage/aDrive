package api

import (
	nn "aDrive/proto/namenode"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DeleteResp struct {
}

func Delete(addr string, absolutePath string) (DeleteResp, error) {
	clientConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	getResp, err := nn.NewNameNodeServiceClient(clientConn).Delete(context.Background(), &nn.DeleteDataReq{
		FileName: absolutePath,
	})
	if err != nil || !getResp.Ok {
		return DeleteResp{}, err
	}
	return DeleteResp{}, err
}
