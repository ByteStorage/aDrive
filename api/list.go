package api

import (
	nn "aDrive/proto/namenode"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ListResp struct {
	Filename []string
	Dirname  []string
}

func List(addr string, path string) (ListResp, error) {
	clientConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	listResp, err := nn.NewNameNodeServiceClient(clientConn).List(context.Background(), &nn.ListReq{
		ParentPath: path,
	})
	if err != nil {
		return ListResp{}, err
	}
	return ListResp{
		Filename: listResp.FileName,
		Dirname:  listResp.DirName,
	}, err
}
