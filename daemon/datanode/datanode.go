package datanode

import (
	"aDrive/datanode"
	"aDrive/pkg/utils"
	dn "aDrive/proto/datanode"
	nn "aDrive/proto/namenode"
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io/ioutil"
	"log"
	"net"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	DefaultDataNodePort = 7000
	DefaultNameNodePort = 9999
)

func Start() {
	nameNodeAddr := "localhost:" + strconv.Itoa(DefaultNameNodePort)
	dataLocation := "data/"
	dataNodeInstance := new(datanode.Server)
	if !strings.HasSuffix(dataLocation, "/") {
		dataLocation = dataLocation + "/"
	}
	_, err := ioutil.ReadDir(dataLocation)
	if err != nil {
		err := os.MkdirAll(dataLocation, 0755)
		if utils.Exit("create dir error", err) {
			return
		}
	}
	dataNodeInstance.DataDirectory = dataLocation
	dataNodeInstance.ServicePort = uint32(DefaultDataNodePort)
	nameNodeHost, nameNodePort, err := net.SplitHostPort(nameNodeAddr)
	if err != nil {
		panic(err)
	}
	portInt, err := strconv.Atoi(nameNodePort)
	if err != nil {
		panic(err)
	}
	dataNodeInstance.NameNodeHost = nameNodeHost
	dataNodeInstance.NameNodePort = uint32(portInt)
	zap.L().Info("Data storage location is " + dataLocation)

	addr := ":" + strconv.Itoa(DefaultDataNodePort)

	listener, err := net.Listen("tcp", addr)
	if utils.Exit("failed to listen", err) {
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
	zap.L().Info("start register to name nodes: " + nameNodeAddr)
	for {
		//向NameNode注册
		if dataNodeInstance.NameNodeHost == "" {
			continue
		}
		conn, err := grpc.Dial(nameNodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
			Addr:       "localhost:" + strconv.Itoa(int(dataNodeInstance.ServicePort)),
			UsedDisk:   usedDisk / 1024 / 1024,
			UsedMem:    usedMem / 1024 / 1024,
			TotalMem:   totalMem / 1024 / 1024,
			TotalDisk:  totalDisk / 1024 / 1024,
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

	zap.L().Info("DataNode daemon started on port: " + strconv.Itoa(DefaultDataNodePort))

	go heartBeatToNameNode(dataNodeInstance, strconv.Itoa(DefaultDataNodePort))

	// graceful shutdown
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig

	server.GracefulStop()
}

func StartServer(host, nameNodeAddr string, serverPort int, dataLocation string) {
	// 开启pprof，监听请求
	go func() {
		router := gin.New()
		pprof.Register(router)
		p := fmt.Sprintf(":%d", 6061)
		if err := router.Run(p); err != nil {
			log.Println(err)
		}
	}()
	fmt.Println("host: ", host)
	fmt.Println("nameNodeAddr: ", nameNodeAddr)
	fmt.Println("serverPort: ", serverPort)
	fmt.Println("dataLocation: ", dataLocation)
	dataNodeInstance := new(datanode.Server)
	if !strings.HasSuffix(dataLocation, "/") {
		dataLocation = dataLocation + "/"
	}
	//判断文件夹是否存在，不存在则创建一个
	_, err := ioutil.ReadDir(dataLocation)
	if err != nil {
		err := os.MkdirAll(dataLocation, 0755)
		if utils.Exit("create dir error", err) {
			return
		}
	}
	dataNodeInstance.DataDirectory = dataLocation
	dataNodeInstance.ServicePort = uint32(serverPort)
	nameNodeHost, nameNodePort, err := net.SplitHostPort(nameNodeAddr)
	if err != nil {
		panic(err)
	}
	portInt, err := strconv.Atoi(nameNodePort)
	if err != nil {
		panic(err)
	}
	fmt.Println("nameNodeHost: ", nameNodeHost)
	fmt.Println("nameNodePort: ", nameNodePort)
	dataNodeInstance.NameNodeHost = nameNodeHost
	dataNodeInstance.NameNodePort = uint32(portInt)

	/*// 获取公网 IP 地址
	ipResp, err := http.Get("https://api.ipify.org")
	if err != nil {
		zap.L().Error("Failed to get public IP address", zap.Error(err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zap.L().Error("Failed to close IP response", zap.Error(err))
		}
	}(ipResp.Body)

	ipBytes, err := ioutil.ReadAll(ipResp.Body)

	if err != nil {
		zap.L().Error("Failed to read IP response", zap.Error(err))
	}
	ipAddress := string(ipBytes)
	apiUrl := fmt.Sprintf("https://ipinfo.io/%s?token=", ipAddress)

	resp, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zap.L().Error("Failed to close IP response", zap.Error(err))
		}
	}(resp.Body)*/

	/*body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var data datanode.AutoGenerated
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	place := data.Country + " " + data.Region + " " + data.City
	dataNodeInstance.Place = place*/

	/*go listenLeader(nameNodeAddr, dataNodeInstance)*/

	zap.L().Info("Data storage location is " + dataLocation)

	addr := ":" + strconv.Itoa(serverPort)

	listener, err := net.Listen("tcp", addr)
	if utils.Exit("failed to listen", err) {
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
	if host == "" || host == "localhost" {
		//为了方便Windows下调试
		host = "localhost"
	}
	zap.L().Info("start register to name nodes: " + nameNodeAddr)
	for true {
		//向NameNode注册
		if dataNodeInstance.NameNodeHost == "" {
			continue
		}
		conn, err := grpc.Dial(nameNodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
			Addr:       host + ":" + strconv.Itoa(int(dataNodeInstance.ServicePort)),
			UsedDisk:   usedDisk / 1024 / 1024,
			UsedMem:    usedMem / 1024 / 1024,
			TotalMem:   totalMem / 1024 / 1024,
			TotalDisk:  totalDisk / 1024 / 1024,
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

	fmt.Println("dataNodeInstance", dataNodeInstance)
	go heartBeatToNameNode(dataNodeInstance, host)

	// graceful shutdown
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig

	server.GracefulStop()

}

func listenLeader(addr string, instance *datanode.Server) {
	/*for range time.Tick(time.Second * 3) {
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
	}*/
}

func heartBeatToNameNode(instance *datanode.Server, host string) {
	i := 0
	for range time.Tick(time.Second * 10) {
		var conn *grpc.ClientConn
		i++
		addr := instance.NameNodeHost + ":" + strconv.Itoa(int(instance.NameNodePort))
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Info("moving to new name node")
		}
		client := nn.NewNameNodeServiceClient(conn)
		_, err = client.HeartBeat(context.Background(), &nn.HeartBeatReq{
			Addr: host + ":" + strconv.Itoa(int(instance.ServicePort)),
		})
		if err != nil {
			zap.L().Info("moving to new name node")
		}
		//每150秒更新信息
		if i >= 30 {
			usedMem, err := datanode.GetUsedMem()
			if err != nil {
				continue
			}
			usedDisk, err := datanode.GetUsedDisk()
			if err != nil {
				continue
			}
			totalMem, err := datanode.GetTotalMem()
			if err != nil {
				continue
			}
			totalDisk, err := datanode.GetTotalDisk()
			if err != nil {
				continue
			}
			cpuPercent, err := datanode.GetCpuPercent()
			if err != nil {
				continue
			}
			resp, err := client.UpdateDataNodeMessage(context.Background(), &nn.UpdateDataNodeMessageReq{
				UsedDisk:   usedDisk / 1024 / 1024,
				UsedMem:    usedMem / 1024 / 1024,
				TotalMem:   totalMem / 1024 / 1024,
				TotalDisk:  totalDisk / 1024 / 1024,
				CpuPercent: float32(cpuPercent),
				Addr:       host + ":" + strconv.Itoa(int(instance.ServicePort)),
				Place:      instance.Place,
			})
			if err != nil {
				continue
			}
			if resp.Success {
				i = 0
			}
		}
		conn.Close()
	}
}
