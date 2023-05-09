package datanode

import (
	dn "aDrive/proto/datanode"
	"context"
	"errors"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func (s *Server) Put(c context.Context, req *dn.PutReq) (*dn.PutResp, error) {
	filename := path.Join(s.DataDirectory, req.AbsolutePath)
	dir := filepath.Dir(filename)
	_, err := ioutil.ReadDir(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			zap.L().Error("create dir error" + err.Error())
			return &dn.PutResp{}, errors.New("create dir error:" + err.Error())
		}

	}
	_, err = os.Create(filename)
	if err != nil {
		zap.L().Error("create file error" + err.Error())
		return &dn.PutResp{}, errors.New("create file error:" + err.Error())
	}
	err = ioutil.WriteFile(filename, req.Data, 0666)
	if err != nil {
		zap.L().Error("write file error" + err.Error())
		return &dn.PutResp{}, errors.New("write file error:" + err.Error())
	}
	return &dn.PutResp{
		Success: true,
	}, nil
}
