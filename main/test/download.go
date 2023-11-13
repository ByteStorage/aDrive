package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

const chunkSize = 1024 * 1024 / 10 // 1MB

func downloadFile(url string, outputPath string, wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()

	// 发送HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Error downloading %s: %s", url, err)
		return
	}
	defer resp.Body.Close()

	// 创建文件来保存下载的内容
	file, err := os.Create(outputPath)
	if err != nil {
		ch <- fmt.Sprintf("Error creating file for %s: %s", url, err)
		return
	}
	defer file.Close()

	// 将HTTP响应的Body拷贝到文件中
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("Error copying content for %s: %s", url, err)
		return
	}

	ch <- fmt.Sprintf("Downloaded %s", url)
}

func downloadFileInChunks(url string, outputPath string, wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()

	// 获取文件大小
	resp, err := http.Head(url)
	if err != nil {
		ch <- fmt.Sprintf("Error getting file information for %s: %s", url, err)
		return
	}
	defer resp.Body.Close()

	fileSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		ch <- fmt.Sprintf("Error parsing Content-Length header for %s: %s", url, err)
		return
	}

	// 使用WaitGroup来等待所有goroutine完成
	var downloadWg sync.WaitGroup

	// 分割文件并启动goroutine下载每个块
	for i := 0; i < fileSize; i += chunkSize {
		downloadWg.Add(1)
		end := i + chunkSize - 1
		if end >= fileSize {
			end = fileSize - 1
		}
		go downloadChunk(url, i, end, &downloadWg, ch)
	}

	// 等待所有下载完成
	downloadWg.Wait()

	// 合并下载的文件块
	partFiles := []string{}
	for i := 0; i < fileSize; i += chunkSize {
		partFiles = append(partFiles, fmt.Sprintf("%s_part%d", outputPath, i/chunkSize))
	}

	err = mergeFiles(outputPath, partFiles)
	if err != nil {
		ch <- fmt.Sprintf("Error merging files for %s: %s", url, err)
		return
	}

	ch <- fmt.Sprintf("Files merged successfully. Output file: %s", outputPath)

	// 删除临时文件块
	for _, partFile := range partFiles {
		err := os.Remove(partFile)
		if err != nil {
			fmt.Printf("Error removing part file %s: %s\n", partFile, err)
		}
	}
}

func downloadChunk(url string, start, end int, wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()

	// 发送HTTP请求，指定Range头部以下载特定范围的数据
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- fmt.Sprintf("Error creating request for %s: %s", url, err)
		return
	}

	rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
	req.Header.Add("Range", rangeHeader)

	resp, err := client.Do(req)
	if err != nil {
		ch <- fmt.Sprintf("Error downloading %s: %s", url, err)
		return
	}
	defer resp.Body.Close()

	// 创建文件来保存下载的内容
	file, err := os.Create(fmt.Sprintf("%s_part%d", getFilename(url), start/chunkSize))
	if err != nil {
		ch <- fmt.Sprintf("Error creating file for %s: %s", url, err)
		return
	}
	defer file.Close()

	// 将HTTP响应的Body拷贝到文件中
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("Error copying content for %s: %s", url, err)
		return
	}

	ch <- fmt.Sprintf("Downloaded %s (bytes %d-%d)", url, start, end)
}

func mergeFiles(outputPath string, partFiles []string) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Error creating output file: %s", err)
	}
	defer outputFile.Close()

	for _, partFile := range partFiles {
		part, err := os.Open(partFile)
		if err != nil {
			return fmt.Errorf("Error opening part file %s: %s", partFile, err)
		}
		defer part.Close()

		_, err = io.Copy(outputFile, part)
		if err != nil {
			return fmt.Errorf("Error copying part file %s: %s", partFile, err)
		}
	}

	return nil
}

func getFilename(url string) string {
	// 使用path包中的Base函数获取URL中的文件名
	return path.Base(url)
}

func main() {
	// 要下载的URL
	start := time.Now()
	url := "https://mirrors.aliyun.com/golang/go1.10.3.linux-amd64.tar.gz?spm=a2c6h.25603864.0.0.32467c45P0r16i"

	// 设置下载后的文件名
	outputPath := getFilename(url)

	// 使用WaitGroup来等待所有goroutine完成
	var wg sync.WaitGroup

	// 使用channel来接收每个goroutine的结果
	resultChannel := make(chan string)

	// 启动一个goroutine来监听结果并打印
	go func() {
		for result := range resultChannel {
			fmt.Println(result)
		}
	}()

	// 启动下载文件的goroutine
	wg.Add(1)
	go downloadFileInChunks(url, outputPath, &wg, resultChannel)

	// 等待下载完成
	wg.Wait()

	// 关闭结果channel
	close(resultChannel)
	fmt.Printf("Total time: %s\n", time.Since(start))
}
