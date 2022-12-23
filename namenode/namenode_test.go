package namenode

import (
	nn "aDrive/proto/namenode"
	"context"
	"testing"
)

func initService() *Service {
	return &Service{
		IdToDataNodes: map[int64]DataNodeInstance{
			1: {Host: "localhost", ServicePort: "50051"},
			2: {Host: "localhost", ServicePort: "50052"},
			3: {Host: "localhost", ServicePort: "50053"},
		},
		FileNameToDataNodes: map[string][]DataMessage{},
		DataNodeMessageMap:  map[string]DataNodeMessage{},
	}
}

func TestService_Put(t *testing.T) {
	service := initService()
	putResp, err := service.Put(context.Background(), &nn.PutReq{
		AbsolutePath: "/test.txt",
		Data:         []byte("hello world"),
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(putResp.AbsolutePath, putResp.DataMessage)
}
