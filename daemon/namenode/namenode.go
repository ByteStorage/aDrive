package namenode

import (
	"aDrive/main/web"
	"aDrive/namenode"
	nn "aDrive/proto/namenode"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/Jille/raftadmin"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	DefaultNameNodePort = 9999
)

var lock sync.Mutex

func Start() {
	addr := ":" + strconv.Itoa(DefaultNameNodePort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		zap.L().Error("listen error", zap.Error(err))
		os.Exit(1)
	}
	nameNodeInstance := namenode.NewService(nil, nil, DefaultNameNodePort)
	server := grpc.NewServer()
	nn.RegisterNameNodeServiceServer(server, nameNodeInstance)
	go func() {
		if err := server.Serve(listener); err != nil {
			zap.L().Error("Server Serve failed in " + addr + ",the reason is " + err.Error())
			os.Exit(1)
		}
	}()
	// graceful shutdown
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig

	server.GracefulStop()
}

func StartServer(host string, master bool, follow string, serverPort int) {
	// 开启pprof，监听请求
	go func() {
		router := gin.New()
		pprof.Register(router)
		p := fmt.Sprintf(":%d", 6060)
		if err := router.Run(p); err != nil {
			log.Println(err)
		}
	}()
	hostname, err := os.Hostname()
	if err != nil {
		zap.L().Error("get hostname error", zap.Error(err))
		os.Exit(1)
	}
	if host == "" || host == "localhost" {
		//为了方便Windows下调试
		hostname = "localhost"
	}
	id := hostname + strconv.Itoa(serverPort)
	baseDir := filepath.Join("./Cluster", id)
	join := filepath.Join(baseDir, "logs.dat")

	zap.L().Info("NameNode port is " + strconv.Itoa(serverPort))
	zap.L().Info("Raft Log Dir is " + join)

	addr := ":" + strconv.Itoa(serverPort)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		zap.L().Error("listen error", zap.Error(err))
		os.Exit(1)
	}

	var fsm namenode.Service
	raftNode, tm, ldb, err := newRaft(baseDir, master, follow, id, hostname+addr, &fsm)
	if err != nil {
		zap.L().Error("start raft cluster fail:" + err.Error())
	}
	nameNodeInstance := namenode.NewService(raftNode, ldb, uint16(serverPort))

	if !master {
		file, err := ioutil.ReadFile("./Cluster/metadata.dat")
		if err != nil {
			zap.L().Error("read metadata.dat error", zap.Error(err))
		}
		var s namenode.Service
		err = json.Unmarshal(file, &s)
		if err != nil {
			zap.L().Error("unmarshal metadata.dat error", zap.Error(err))
		}
		if s.DirTree != nil {
			nameNodeInstance.DirTree = s.DirTree
			nameNodeInstance.IdToDataNodes = s.IdToDataNodes
			nameNodeInstance.FileNameToDataNodes = s.FileNameToDataNodes
			nameNodeInstance.DataNodeHeartBeat = s.DataNodeHeartBeat
			nameNodeInstance.DataNodeMessageMap = s.DataNodeMessageMap
		}
	}

	// 注册prometheus
	// Create a metrics registry.
	prometheusReg := prometheus.NewRegistry()
	// Create some standard server metrics.
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	// Register standard server metrics to registry.
	prometheusReg.MustRegister(grpcMetrics)

	server := grpc.NewServer()
	nn.RegisterNameNodeServiceServer(server, nameNodeInstance)
	tm.Register(server)
	raftadmin.Register(server, raftNode)
	reflection.Register(server)
	grpcMetrics.InitializeMetrics(server)

	go func() {
		if err := server.Serve(listener); err != nil {
			zap.L().Error("Server Serve failed in " + addr + ",the reason is " + err.Error())
			os.Exit(1)
		}
	}()

	zap.L().Info("NameNode daemon started on port: " + strconv.Itoa(serverPort))

	go func(s *namenode.Service) {
		for range time.Tick(30 * time.Second) {
			if s.RaftNode == nil {
				log.Println("节点加入有误，raft未建立成功")
				continue
			}
			zap.L().Debug("", zap.Any("idToDataNodes", s.IdToDataNodes))
		}
	}(nameNodeInstance)

	listenPath := filepath.Join("./Cluster", "metadata.dat")
	_, err = os.Create(listenPath)
	if err != nil {
	}

	go listenLeaderChanges(listenPath, nameNodeInstance)

	// 监测datanode的心跳
	go checkDataNode(nameNodeInstance)

	go sendMessageToProm(nameNodeInstance)

	go web.StartWeb()

	// graceful shutdown
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig

	server.GracefulStop()

}

func sendMessageToProm(s *namenode.Service) {
	for range time.Tick(15 * time.Second) {
		fmt.Println("sendMessageToProm方法执行，此时goruoutine数量为：", runtime.NumGoroutine())
		//将map转换为数组
		var dataNodeMessages []namenode.DataNodeMessage
		for _, v := range s.DataNodeMessageMap {
			dataNodeMessages = append(dataNodeMessages, v)
			web.UsedMem.Set(float64(v.UsedMem))
			web.FreeMem.Set(float64(v.TotalMem - v.UsedMem))
			web.UsedDisk.Set(float64(v.UsedDisk))
			web.FreeDisk.Set(float64(v.TotalDisk - v.UsedDisk))
		}
		for i, v := range dataNodeMessages {
			web.Place.With(prometheus.Labels{"place": "节点" + strconv.Itoa(i+1) + ":      " + v.Place}).Set(float64(i + 1))
		}
	}
}

// 监听主节点元数据信息变化
func listenLeaderChanges(filename string, s *namenode.Service) {
	/*watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func(s *namenode.Service) {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if s.RaftNode.State() == raft.Leader {
						continue
					}

					bytes, err := ioutil.ReadFile(filename)
					if err != nil {
						zap.L().Error("cannot read file:" + err.Error())
					}

					var srv namenode.Service
					err = json.Unmarshal(bytes, &srv)
					if err != nil {
						zap.L().Error("cannot parse json:" + err.Error())
					}
					copyService(s, srv)
					address, _ := s.RaftNode.LeaderWithID()
					web.CurrentLeader = string(address)
					continue
				}
			case err := <-watcher.Errors:
				zap.L().Error("watch error:" + err.Error())
				continue
			}
		}
	}(s)
	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}
	<-done*/
}

func checkDataNode(instance *namenode.Service) {
	for range time.Tick(time.Second * 10) {
		fmt.Println("checkDataNode方法执行，此时goruoutine数量为：", runtime.NumGoroutine())
		if instance.RaftNode.State() != raft.Leader {
			continue
		}
		if instance.IdToDataNodes == nil {
			continue
		}
		for k, v := range instance.IdToDataNodes {
			addr := v.Host + ":" + v.ServicePort
			lastHeartBeatTime := instance.DataNodeHeartBeat[addr]
			if lastHeartBeatTime.Add(time.Second*15).Unix() < time.Now().Unix() {
				/*var reply bool
				err := instance.ReDistributeData(&namenode.ReDistributeDataRequest{DataNodeUri: addr}, &reply)
				if err != nil {
					zap.L().Error("cannot redistribute data:" + err.Error())
				}*/
				lock.Lock()
				delete(instance.IdToDataNodes, k)
				lock.Unlock()
			}
		}
	}
}

func newRaft(baseDir string, master bool, follow, myID, myAddress string, fsm raft.FSM) (*raft.Raft, *transport.Manager, *boltdb.BoltStore, error) {
	c := raft.DefaultConfig()
	isLeader := make(chan bool, 1)
	c.NotifyCh = isLeader
	c.LocalID = raft.ServerID(myID)

	if master {
		err := os.RemoveAll("./Cluster")
		if err != nil {
		}
	}
	err := os.MkdirAll(baseDir, 0755)
	if err != nil {
		return nil, nil, nil, err
	}
	ldb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "logs.dat"), err)
	}

	sdb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "stable.dat"), err)
	}

	fss, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(`raft.NewFileSnapshotStore(%q, ...): %v`, baseDir, err)
	}

	tm := transport.New(raft.ServerAddress(myAddress), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})

	r, err := raft.NewRaft(c, fsm, ldb, sdb, fss, tm.Transport())
	if err != nil {
		zap.L().Error("cannot create raft:" + err.Error())
		return nil, nil, nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	if master {
		zap.L().Info("start bootstrap cluster")

		web.CurrentLeader = myAddress
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(myID),
					Address:  raft.ServerAddress(myAddress),
				},
			},
		}
		f := r.BootstrapCluster(cfg)
		if err := f.Error(); err != nil {
			zap.L().Error("bootstrap cluster error:" + err.Error())
			return nil, nil, nil, fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
		zap.L().Info("bootstrap cluster success,then start web server")
	} else if follow != "" {
		findLeader, err := FindLeader(follow)
		if err != nil {
			zap.L().Error("cannot find leader:" + err.Error())
			return nil, nil, nil, err
		}
		web.CurrentLeader = findLeader
		conn, err := grpc.Dial(findLeader, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Error("cannot connect to leader:" + err.Error())
			return nil, nil, nil, err
		}
		resp, err := nn.NewNameNodeServiceClient(conn).JoinCluster(context.Background(), &nn.JoinClusterReq{
			Addr:          myAddress,
			Id:            myID,
			PreviousIndex: 0,
		})
		if err != nil {
			zap.L().Error("cannot join cluster:" + err.Error())
			return nil, nil, nil, err
		}
		if resp.Success {
			log.Println("join the cluster success")
		}

	}

	return r, tm, ldb, nil
}

// FindLeader 查找NameNode的Raft集群的Leader
func FindLeader(addrList string) (string, error) {

	zap.L().Debug(addrList)
	nameNodes := strings.Split(addrList, ",")
	var res = ""
	var conn *grpc.ClientConn
	defer conn.Close()
	var err error
	for _, n := range nameNodes {
		conn, err = grpc.Dial(n, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Error("cannot connect to:" + n)
			//表明连接不上，继续	遍历节点
			continue
		}
		resp, err := nn.NewNameNodeServiceClient(conn).FindLeader(context.Background(), &nn.FindLeaderReq{})
		if err != nil {
			zap.L().Error("cannot find leader:" + err.Error())
			continue
		}
		res = resp.Addr
		break
	}
	if res == "" {
		err = errors.New("there is no alive name node")
		if err != nil {
			return "", err
		}
	}
	return res, nil
}

func copyService(old *namenode.Service, new namenode.Service) {
	old.FileNameToDataNodes = new.FileNameToDataNodes
	old.DirTree = new.DirTree
	old.IdToDataNodes = new.IdToDataNodes
	old.DataNodeHeartBeat = new.DataNodeHeartBeat
	old.DataNodeMessageMap = new.DataNodeMessageMap
}
