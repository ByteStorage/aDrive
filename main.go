package main

import (
	"aDrive/main/cli"
	"aDrive/pkg/logger"
	_ "net/http/pprof"
	"sync"
)

var once sync.Once

func init() {
	once.Do(func() {
		logger.Init()
	})
}

func main() {

	cli.StartServer()
}
