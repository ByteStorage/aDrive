package api

import (
	"aDrive/daemon/namenode"
	nn "aDrive/proto/namenode"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var ApiService = &Service{}

type Service struct {
	NameNodeHost string
	NameNodePort string
}

func InitApi(nameNodeAddress string) (*grpc.ClientConn, error) {
	address := strings.Split(nameNodeAddress, ",")
	if len(address) >= 3 {
		go listenLeader(ApiService, nameNodeAddress)
	} else {
		ApiService.NameNodePort = strconv.Itoa(namenode.DefaultNameNodePort)
		ApiService.NameNodeHost = address[0]
	}
	for {
		if ApiService.NameNodeHost != "" {
			break
		}
	}
	return grpc.Dial(ApiService.NameNodeHost+":"+ApiService.NameNodePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func listenLeader(s *Service, address string) {
	for range time.Tick(time.Second * 1) {
		log.Println(s.NameNodeHost, s.NameNodePort)
		log.Println(address)
		nameNodes := strings.Split(address, ",")
		for _, n := range nameNodes {
			conn, err := grpc.Dial(n, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				//表明连接不上，继续遍历节点
				log.Println(err)
				continue
			}
			resp, err := nn.NewNameNodeServiceClient(conn).FindLeader(context.Background(), &nn.FindLeaderReq{})
			if err != nil {
				log.Println(err)
				continue
			}
			host, port, err := net.SplitHostPort(resp.Addr)
			if err != nil {
				panic(err)
			}
			s.NameNodeHost = host
			s.NameNodePort = port
		}
		log.Println(s.NameNodeHost, s.NameNodePort)
		if s.NameNodePort == "" {
			err := errors.New("there is no alive name node")
			if err != nil {
				panic(err)
			}
		}
	}
}
