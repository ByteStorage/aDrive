package datanode

import (
	dn "aDrive/proto/datanode"
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
)

func (s *Server) Stat(c context.Context, req *dn.StatReq) (*dn.StatResp, error) {
	p := path.Join(s.DataDirectory + req.PrePath + req.BlockId)
	zap.S().Debug("Stat: ", p)
	stat, err := os.Stat(p)
	if err != nil {
		zap.S().Error("cannot stat the file:", err)
		return &dn.StatResp{}, err
	}
	_, name := filepath.Split(stat.Name())
	mod := fmt.Sprintf("%s", stat.Mode())
	return &dn.StatResp{
		Size:    stat.Size(),
		ModTime: stat.ModTime().Unix(),
		Name:    name,
		Mode:    mod,
	}, nil
}
