package namenode

import (
	"aDrive/pkg/utils"
	dn "aDrive/proto/datanode"
	nn "aDrive/proto/namenode"
	"context"
	"fmt"
	"github.com/klauspost/reedsolomon"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *Service) Get(c context.Context, req *nn.GetReq) (*nn.GetResp, error) {
	absolutePath := req.AbsolutePath
	ext := filepath.Ext(absolutePath)
	prefix := strings.TrimSuffix(absolutePath, ext)
	dataMessages := s.FileNameToDataNodes[req.AbsolutePath]
	fmt.Println(s.FileNameToDataNodes)
	data := make([][]byte, len(dataMessages))
	for i := range dataMessages {
		//向各个DataNode发送Get请求
		clientConn, err := grpc.Dial(dataMessages[i].Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			data[dataMessages[i].Id] = nil
			continue
		}
		client := dn.NewDataNodeClient(clientConn)
		getResp, err := client.Get(context.Background(), &dn.GetReq{
			AbsolutePath: prefix + strconv.Itoa(dataMessages[i].Id) + ext,
		})
		if err != nil {
			data[dataMessages[i].Id] = nil
			continue
		}
		data[dataMessages[i].Id] = getResp.Data
	}
	log.Println("dataMessages", dataMessages)
	encoder, err := reedsolomon.New(len(dataMessages)-len(dataMessages)/3, len(dataMessages)/3)
	if utils.Exit("new encoder error", err) {
		return &nn.GetResp{}, err
	}
	verify, err := encoder.Verify(data)
	if !verify {
		err := encoder.Reconstruct(data)
		if utils.Exit("reconstruct error", err) {
			return &nn.GetResp{}, err
		}
		b, err := encoder.Verify(data)
		if utils.Exit("verify error", err) {
			return &nn.GetResp{}, err
		}
		if !b {
			zap.L().Error("verify error")
			return &nn.GetResp{}, err
		}
	}
	var res []byte
	for i := range data {
		res = append(res, data[i]...)
	}
	return &nn.GetResp{
		Data: res,
	}, nil
}
