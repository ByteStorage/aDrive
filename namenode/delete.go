package namenode

import (
	"aDrive/pkg/utils"
	dn "aDrive/proto/datanode"
	nn "aDrive/proto/namenode"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (s *Service) Delete(c context.Context, req *nn.DeleteDataReq) (*nn.DeleteDataResp, error) {
	dirTreeFilePath := utils.ModPath(req.FileName)

	ext := filepath.Ext(req.FileName)
	prefix := strings.TrimSuffix(req.FileName, ext)

	var res nn.DeleteDataResp
	s.DirTree.Delete(s.DirTree.Root, dirTreeFilePath)
	zap.S().Debug("删除文件后目录树为:", s.DirTree.LookAll())

	fileBlocks := s.FileNameToDataNodes[req.FileName]

	for i := range fileBlocks {
		conn, err := grpc.Dial(fileBlocks[i].Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Error("connect to datanode(" + fileBlocks[i].Host + "):" + err.Error())
			return &nn.DeleteDataResp{}, err
		}
		client := dn.NewDataNodeClient(conn)
		_, err = client.Delete(context.Background(), &dn.DeleteReq{
			Filename: prefix + strconv.Itoa(fileBlocks[i].Id) + ext,
		})
		if err != nil {
			zap.L().Error("delete file(" + prefix + strconv.Itoa(fileBlocks[i].Id) + ext + ") error:" + err.Error())
			return &nn.DeleteDataResp{}, err
		}
	}
	delete(s.FileNameToDataNodes, req.FileName)
	zap.S().Debug("删除文件后FileNameToDataNodes为:", s.FileNameToDataNodes)

	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println("cannot marshal data")
		return &nn.DeleteDataResp{}, err
	}
	if s.RaftNode != nil {
		s.RaftNode.Apply(bytes, time.Second*1)
	}
	return &res, nil
}
