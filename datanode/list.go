package datanode

import (
	dn "aDrive/proto/datanode"
	"context"
	"log"
	"os"
)

// List 索引目录下的所有文件和文件夹
func (s *Server) List(c context.Context, req *dn.ListReq) (*dn.ListResp, error) {
	files, err := os.ReadDir(req.Path)
	if err != nil {
		log.Println("cannot list the file:", err)
		return &dn.ListResp{}, err
	}
	var blocks, dirs []string
	for _, file := range files {
		// 判断是否为文件夹
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		} else {
			blocks = append(blocks, file.Name())
		}
	}
	return &dn.ListResp{FileList: blocks, DirList: dirs}, nil
}
