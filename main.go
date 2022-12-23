package main

import (
	"aDrive/pkg/logger"
	"go.uber.org/zap"
	"sync"
)

var once sync.Once

func init() {
	once.Do(func() {
		logger.Init()
	})
}

func main() {
	zap.L().Debug("debug")
	zap.L().Info("info")
	zap.L().Warn("warn")
	zap.L().Error("error")
	zap.S().Debug("Debug")
	zap.S().Info("Info")
	zap.S().Warn("Warn")
	zap.S().Error("Error")
}
