package datanode

import (
	dn "aDrive/proto/datanode"
	"context"
	"errors"
	"go.uber.org/zap"
	"io/ioutil"
)

func (s *Server) Put(c context.Context, req *dn.PutReq) (*dn.PutResp, error) {
	err := ioutil.WriteFile(req.AbsolutePath, req.Data, 0666)
	if err != nil {
		zap.L().Error("write file error" + err.Error())
		return &dn.PutResp{}, errors.New("write file error:" + err.Error())
	}
	return &dn.PutResp{
		Success: true,
	}, nil
}
