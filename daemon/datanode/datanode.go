package datanode

import (
	"aDrive/datanode"
	dn "aDrive/proto/datanode"
	nn "aDrive/proto/namenode"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func StartServer(host, nameNodeAddr string, serverPort int, dataLocation string) {
	dataNodeInstance := new(datanode.Server)
	if !strings.HasSuffix(dataLocation, "/") {
		dataLocation = dataLocation + "/"
	}
	dataNodeInstance.DataDirectory = dataLocation
	dataNodeInstance.ServicePort = uint32(serverPort)

	go listenLeader(nameNodeAddr, dataNodeInstance)

	zap.L().Info("Data storage location is " + dataLocation)

	addr := ":" + strconv.Itoa(serverPort)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		zap.L().Error("failed to listen: " + err.Error())
		return
	}

	server := grpc.NewServer()
	dn.RegisterDataNodeServer(server, dataNodeInstance)

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Printf(fmt.Sprintf("Server Serve failed in %s", addr), "err", err.Error())
			panic(err)
		}
	}()

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	if host == "" {
		//为了方便Windows下调试
		hostname = "localhost"
	}
	zap.L().Info("start register to name nodes: " + nameNodeAddr)
	for true {
		//向NameNode注册
		if dataNodeInstance.NameNodeHost == "" {
			continue
		}
		addr := dataNodeInstance.NameNodeHost + ":" + strconv.Itoa(int(dataNodeInstance.NameNodePort))
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		usedMem, err := datanode.GetUsedMem()
		if err != nil {

		}
		usedDisk, err := datanode.GetUsedDisk()
		if err != nil {

		}
		totalMem, err := datanode.GetTotalMem()
		if err != nil {

		}
		totalDisk, err := datanode.GetTotalDisk()
		if err != nil {

		}
		cpuPercent, err := datanode.GetCpuPercent()
		if err != nil {

		}
		resp, err := nn.NewNameNodeServiceClient(conn).RegisterDataNode(context.Background(), &nn.RegisterDataNodeReq{
			Addr:       hostname + ":" + strconv.Itoa(int(dataNodeInstance.ServicePort)),
			UsedDisk:   usedDisk,
			UsedMem:    usedMem,
			TotalMem:   totalMem,
			TotalDisk:  totalDisk,
			CpuPercent: float32(cpuPercent),
		})
		if err != nil {
			continue
		}
		if resp.Success {
			zap.L().Info("register success")
			break
		}
	}

	zap.L().Info("DataNode daemon started on port: " + strconv.Itoa(serverPort))

	go heartBeatToNameNode(dataNodeInstance, hostname)

	// graceful shutdown
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig

	server.GracefulStop()

}

func listenLeader(addr string, instance *datanode.Server) {
	for range time.Tick(time.Second * 1) {
		nameNodes := strings.Split(addr, ",")
		for _, n := range nameNodes {
			conn, err := grpc.Dial(n, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				//表明连接不上，继续遍历节点
				zap.L().Info("connect to name node failed: " + err.Error())
				continue
			}
			resp, err := nn.NewNameNodeServiceClient(conn).FindLeader(context.Background(), &nn.FindLeaderReq{})
			if err != nil {
				zap.L().Info("find leader failed: " + err.Error())
				continue
			}
			host, port, err := net.SplitHostPort(resp.Addr)
			if err != nil {
				zap.L().Error("split host port failed: " + err.Error())
				os.Exit(1)
			}
			instance.NameNodeHost = host
			p, err := strconv.Atoi(port)
			if err != nil {
				zap.L().Error("parse port failed: " + err.Error())
				os.Exit(1)
			}
			instance.NameNodePort = uint32(p)
			break
		}
		if instance.NameNodeHost == "" {
			err := errors.New("there is no alive name node")
			if err != nil {
				zap.L().Error("there is no alive name node")
				os.Exit(1)
			}
		}
		zap.L().Info("DataNode信息为:" + string(instance.NameNodeHost) + ":" + strconv.Itoa(int(instance.NameNodePort)))
	}
}

func heartBeatToNameNode(instance *datanode.Server, hostname string) {
	i := 0
	for range time.Tick(time.Second * 10) {
		i++
		addr := instance.NameNodeHost + ":" + strconv.Itoa(int(instance.NameNodePort))
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Info("moving to new name node")
		}
		client := nn.NewNameNodeServiceClient(conn)
		_, err = client.HeartBeat(context.Background(), &nn.HeartBeatReq{
			Addr: hostname + ":" + strconv.Itoa(int(instance.ServicePort)),
		})
		if err != nil {
			zap.L().Info("moving to new name node")
		}
		//每150秒更新信息
		if i >= 30 {
			usedMem, err := datanode.GetUsedMem()
			if err != nil {

			}
			usedDisk, err := datanode.GetUsedDisk()
			if err != nil {

			}
			totalMem, err := datanode.GetTotalMem()
			if err != nil {

			}
			totalDisk, err := datanode.GetTotalDisk()
			if err != nil {

			}
			cpuPercent, err := datanode.GetCpuPercent()
			if err != nil {

			}
			resp, err := client.UpdateDataNodeMessage(context.Background(), &nn.UpdateDataNodeMessageReq{
				UsedDisk:   usedDisk,
				UsedMem:    usedMem,
				TotalMem:   totalMem,
				TotalDisk:  totalDisk,
				CpuPercent: float32(cpuPercent),
				Addr:       hostname + ":" + strconv.Itoa(int(instance.ServicePort)),
			})
			if err != nil {

			}
			if resp.Success {
				i = 0
			}
		}
	}
}
