package datanode

import (
	dn "aDrive/proto/datanode"
	"context"
	"log"
	"os"
)

// Rename 重命名文件
func (s *Server) Rename(c context.Context, req *dn.RenameReq) (*dn.RenameResp, error) {
	err := os.Rename(s.DataDirectory+req.OldPath, s.DataDirectory+req.NewPath)
	if err != nil {
		log.Println("cannot rename the file:", err)
		return &dn.RenameResp{}, err
	}
	log.Println("成功重命名")
	return &dn.RenameResp{Success: true}, nil
}
