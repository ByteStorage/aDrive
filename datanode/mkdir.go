package datanode

import (
	dn "aDrive/proto/datanode"
	"context"
	"go.uber.org/zap"
	"os"
	"strings"
)

// Mkdir 创建目录
func (s *Server) Mkdir(c context.Context, req *dn.MkdirReq) (*dn.MkdirResp, error) {
	//判断用户是否携带/
	path := req.Path
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	err := os.MkdirAll(s.DataDirectory+path, 0755)
	if err != nil {
		zap.S().Error("cannot mkdir the file:", err)
		return &dn.MkdirResp{}, err
	}
	zap.S().Info("成功创建目录")
	return &dn.MkdirResp{Success: true}, nil
}
