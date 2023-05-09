package datanode

import (
	dn "aDrive/proto/datanode"
	"context"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
)

func (s *Server) Get(c context.Context, req *dn.GetReq) (*dn.GetResp, error) {
	filename := path.Join(s.DataDirectory, req.AbsolutePath)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		zap.L().Error("read file error" + err.Error())
		return &dn.GetResp{}, err
	}
	return &dn.GetResp{
		Data: bytes,
	}, nil
}
