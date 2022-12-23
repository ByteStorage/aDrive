package datanode

import (
	dn "aDrive/proto/datanode"
	"context"
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"log"
	"os"
	"time"
)

type Server struct {
	dn.DataNodeServer
	DataDirectory string
	ServicePort   uint32
	NameNodeHost  string
	NameNodePort  uint32
}

func (s *Server) Ping(c context.Context, req *dn.PingReq) (*dn.PingResp, error) {
	//接收到NameNode的Ping请求
	s.NameNodeHost = req.Host
	s.NameNodePort = req.Port
	log.Println("I am alive")
	return &dn.PingResp{Success: true}, nil
}

func (s *Server) HeartBeat(c context.Context, req *dn.HeartBeatReq) (*dn.HeartBeatResp, error) {
	if req.Request {
		log.Println("receive heart beat success")
		//以下可改成协程进行，作用不大，heartbeat已经是五秒一请求
		diskPercent, err := GetUsedDisk()
		if err != nil {
			log.Println("cannot GetUsedDisk:", err)
		}
		memPercent, err := GetUsedMem()
		if err != nil {
			log.Println("cannot GetUsedMem:", err)
		}
		cpuPercent, err := GetCpuPercent()
		if err != nil {
			log.Println("cannot GetCpuPercent:", err)
		}
		totalDisk, err := GetTotalDisk()
		if err != nil {
			log.Println("cannot GetTotalDisk:", err)
		}
		totalMem, err := GetTotalMem()
		if err != nil {
			log.Println("cannot GetTotalMem:", err)
		}
		return &dn.HeartBeatResp{
			Success:    true,
			UsedDisk:   diskPercent,
			UsedMem:    memPercent,
			CpuPercent: float32(cpuPercent),
			TotalDisk:  totalDisk,
			TotalMem:   totalMem,
		}, nil
	}
	return nil, errors.New("HeartBeatError")
}

// GetCpuPercent 以下方法可以用于给NameNode决定选取哪一个datanode作为写入节点，已测试过，和linux命令行输出的结果相差无几
// GetCpuPercent 获取CPU使用率
func GetCpuPercent() (float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println("Cannot Read CPU Message:", err)
		return 0, err
	}
	return percent[0], nil
}

// GetUsedMem 获取内存已经使用量
func GetUsedMem() (uint64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Cannot Get Memory Percent:", err)
		return 0, err
	}
	return memInfo.Used, nil
}

// GetUsedDisk 获取当前程序所在目录的硬盘已使用字节数量
func GetUsedDisk() (uint64, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("cannot get pwd:", err)
		return 0, err
	}
	usage, err := disk.Usage(pwd)
	if err != nil {
		log.Println("Cannot Usage Disk Usage:", err)
		return 0, err
	}
	return usage.Used, nil
}

// GetTotalMem 获取总内存，方便计算内存占用率
func GetTotalMem() (uint64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Cannot Get Memory Percent:", err)
		return 0, err
	}
	return memInfo.Total, nil
}

// GetTotalDisk 获取磁盘大小，方便计算磁盘利用率
func GetTotalDisk() (uint64, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("cannot get pwd:", err)
		return 0, err
	}
	usage, err := disk.Usage(pwd)
	if err != nil {
		log.Println("Cannot Usage Disk Usage:", err)
		return 0, err
	}
	return usage.Total, nil
}
