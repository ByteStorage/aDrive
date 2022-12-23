package namenode

import (
	"aDrive/namenode"
	nn "aDrive/proto/namenode"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/Jille/raftadmin"
	"github.com/fsnotify/fsnotify"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func StartServer(host string, master bool, follow, raftId string, serverPort int, blockSize int, replicationFactor int) {

	hostname, err := os.Hostname()
	if err != nil {
		zap.L().Error("get hostname error", zap.Error(err))
		os.Exit(1)
	}
	if host == "" {
		//为了方便Windows下调试
		hostname = "localhost"
	}
	baseDir := filepath.Join("./Cluster", raftId)
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
	raftNode, tm, ldb, err := newRaft(baseDir, master, follow, raftId, hostname+addr, &fsm)
	if err != nil {
		zap.L().Info("start raft cluster fail:" + err.Error())
	}
	nameNodeInstance := namenode.NewService(raftNode, ldb, uint64(blockSize), uint64(replicationFactor), uint16(serverPort))

	// 注册prometheus
	// Create a metrics registry.
	prometheusReg := prometheus.NewRegistry()
	// Create some standard server metrics.
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	// Register standard server metrics to registry.
	prometheusReg.MustRegister(grpcMetrics)

	server := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)

	nn.RegisterNameNodeServiceServer(server, nameNodeInstance)
	tm.Register(server)
	raftadmin.Register(server, raftNode)
	reflection.Register(server)

	// Initialize all metrics.
	grpcMetrics.InitializeMetrics(server)

	// TODO: HARD CODING
	p := 9092
	// Start your http server for prometheus.
	go func() {
		initErr := errors.New("init")

		for initErr != nil {
			// Create a HTTP server for prometheus.
			httpServer := &http.Server{
				Handler: promhttp.HandlerFor(
					prometheusReg,
					promhttp.HandlerOpts{},
				),
				// TODO: HARD CODING
				Addr: ":" + strconv.Itoa(p),
			}

			p += 1

			initErr = httpServer.ListenAndServe()
		}

	}()

	zap.L().Info("server for prometheus started on port: " + strconv.Itoa(p-1))

	go func() {
		if err := server.Serve(listener); err != nil {
			zap.L().Info("Server Serve failed in " + addr + ",the reason is " + err.Error())
			os.Exit(1)
		}
	}()

	zap.L().Info("NameNode daemon started on port: " + strconv.Itoa(serverPort))

	go func(s *namenode.Service) {
		for range time.Tick(10 * time.Second) {
			if s.RaftNode == nil {
				log.Println("节点加入有误，raft未建立成功")
				continue
			}
			id, serverId := raftNode.LeaderWithID()
			log.Println("当前主节点名称:", id, serverId)
			log.Println("当前节点的元数据IdToDataNodes信息:", s.IdToDataNodes)
			log.Println("当前节点的元数据DataNodeHeartBeat信息:", s.DataNodeHeartBeat)
		}
	}(nameNodeInstance)

	listenPath := filepath.Join("./Cluster", "metadata.dat")
	_, err = os.Create(listenPath)
	if err != nil {
	}

	go listenLeaderChanges(listenPath, nameNodeInstance)

	// 监测datanode的心跳
	go checkDataNode(nameNodeInstance)

	// graceful shutdown
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig

	server.GracefulStop()

}

// 监听主节点元数据信息变化
func listenLeaderChanges(filename string, s *namenode.Service) {
	watcher, err := fsnotify.NewWatcher()
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
	<-done
}

func checkDataNode(instance *namenode.Service) {
	for range time.Tick(time.Millisecond * 3) {
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
				delete(instance.IdToDataNodes, k)
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
		return nil, nil, nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	log.Println("master", master)
	if master {
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
			return nil, nil, nil, fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
	} else if follow != "" {
		findLeader, err := FindLeader(follow)
		if err != nil {
			return nil, nil, nil, err
		}
		conn, err := grpc.Dial(findLeader, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, nil, nil, err
		}
		resp, err := nn.NewNameNodeServiceClient(conn).JoinCluster(context.Background(), &nn.JoinClusterReq{
			Addr:          myAddress,
			Id:            myID,
			PreviousIndex: 0,
		})
		if err != nil {
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
	nameNodes := strings.Split(addrList, ",")
	var res = ""
	for _, n := range nameNodes {
		conn, err := grpc.Dial(n, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			//表明连接不上，继续遍历节点
			continue
		}
		resp, err := nn.NewNameNodeServiceClient(conn).FindLeader(context.Background(), &nn.FindLeaderReq{})
		if err != nil {
			continue
		}
		res = resp.Addr
		break
	}
	if res == "" {
		err := errors.New("there is no alive name node")
		if err != nil {
			return "", err
		}
	}
	return res, nil
}

func copyService(old *namenode.Service, new namenode.Service) {
	old.FileNameToBlocks = new.FileNameToBlocks
	old.DirTree = new.DirTree
	old.IdToDataNodes = new.IdToDataNodes
	old.DataNodeHeartBeat = new.DataNodeHeartBeat
	old.DataNodeMessageMap = new.DataNodeMessageMap
	old.BlockToDataNodeIds = new.BlockToDataNodeIds
}
