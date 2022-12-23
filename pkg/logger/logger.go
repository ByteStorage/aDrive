package logger

import (
	"aDrive/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
)

var logger *zap.Logger

const (
	ENV_LOG_LEVEL      = "zap_level"
	LOG_LOCATION       = "./logs/runtime.log"
	ERROR_LOG_LOCATION = "./logs/runtime_err.log"
)

func init() {
	locations := []string{LOG_LOCATION, ERROR_LOG_LOCATION}
	for _, location := range locations {
		filePathExist, err := utils.PathExist(location)
		if err != nil {
			panic(err)
		}
		if !filePathExist {
			dir, _ := filepath.Split(location)
			err := os.MkdirAll(dir, 0750)
			if err != nil {
				panic(err)
			}
			_, err = os.Create(location)
			if err != nil {
				panic(err)
			}
		}
	}

}

func Init() {
	encoder := getEncoder()

	// 默认DebugLevel
	level := zapcore.DebugLevel
	l := os.Getenv(ENV_LOG_LEVEL)
	var c1, c2 zapcore.Core
	if l == "release" {
		level = zapcore.InfoLevel
		// FILE ONLY
		c1 = zapcore.NewCore(encoder, getFileWriter(LOG_LOCATION), level)
		c2 = zapcore.NewCore(encoder, getFileWriter(ERROR_LOG_LOCATION), zap.ErrorLevel)
	} else {
		// STD ONLY
		c1 = zapcore.NewCore(encoder, getWriter(), level)
		// STD and FILE
		c2 = zapcore.NewCore(encoder, getWriterWithFile(ERROR_LOG_LOCATION), zap.ErrorLevel)
	}

	core := zapcore.NewTee(c1, c2)
	logger = zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriterWithFile STD and FILE
func getWriterWithFile(location string) zapcore.WriteSyncer {
	file, err := os.OpenFile(location, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	mr := io.MultiWriter(file, os.Stdout)
	return zapcore.AddSync(mr)
}

// getLogWriter STD only
func getWriter() zapcore.WriteSyncer {
	writer := io.Writer(os.Stdout)
	return zapcore.AddSync(writer)
}

// getFileLogWriter FILE only
func getFileWriter(location string) zapcore.WriteSyncer {
	file, err := os.OpenFile(location, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	writer := io.Writer(file)

	return zapcore.AddSync(writer)
}
