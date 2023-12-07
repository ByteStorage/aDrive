package namenode

import (
	"aDrive/pkg/utils"
	nn "aDrive/proto/namenode"
	"context"
	"encoding/json"
	"log"
	"time"
)

func (s *Service) List(c context.Context, req *nn.ListReq) (*nn.ListResp, error) {
	path := utils.ModPath(req.ParentPath)
	dir := s.DirTree.FindSubDir(path)
	var dirNameList []string
	var fileNameList []string
	for _, str := range dir {
		if str == ".." {
			dirNameList = append(dirNameList, str)
			continue
		}
		resp, err := s.IsDir(context.Background(), &nn.IsDirReq{
			Filename: path + str + "/",
		})
		if err != nil {
			log.Println("NameNode IsDir Error:", err)
			return &nn.ListResp{}, err
		}
		if resp.Ok {
			//是目录
			dirNameList = append(dirNameList, str)
		} else {
			fileNameList = append(fileNameList, str)
		}
	}
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println("cannot marshal data")
		return &nn.ListResp{}, err
	}
	if s.RaftNode != nil {
		s.RaftNode.Apply(bytes, time.Second*1)
	}
	return &nn.ListResp{
		FileName: fileNameList,
		DirName:  dirNameList,
	}, nil

}
