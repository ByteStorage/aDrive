package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ModPath 修改path格式
func ModPath(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// ModFilePath 修改path格式
func ModFilePath(path string) string {
	if strings.HasSuffix(path, "/") {
		path = strings.TrimRight(path, "/")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// GetPrePath 获取文件名的前缀路径
func GetPrePath(filename string) string {
	dir, _ := filepath.Split(filename)
	return dir
}
