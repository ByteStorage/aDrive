package datanode

import (
	dn "aDrive/proto/datanode"
	"context"
	"go.uber.org/zap"
	"os"
)

func (s *Server) Delete(c context.Context, req *dn.DeleteReq) (*dn.DeleteResp, error) {
	_, err := os.Open(s.DataDirectory + req.PrePath + req.BlockId)
	zap.S().Debug("will open: "+s.DataDirectory+req.PrePath+req.BlockId, " required: ", req.PrePath, req.BlockId)
	if err != nil {
		zap.S().Debug("文件已经被删掉")
		return &dn.DeleteResp{Success: true}, nil
	}
	err = os.Remove(s.DataDirectory + req.PrePath + req.BlockId)
	if err != nil {
		return &dn.DeleteResp{}, err
	}
	zap.S().Debug("成功删除文件")
	return &dn.DeleteResp{Success: true}, nil
}
