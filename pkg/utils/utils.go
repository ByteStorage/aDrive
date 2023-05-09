package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func Exit(msg string, err error) bool {
	if err != nil {
		zap.L().Error(msg+":  ", zap.Error(err))
		return true
	}
	return false
}

func CreateToken(username string) (string, error) {
	//创建JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte("iDfjwuhfDasjdnJBDhfSDas"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(token string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("iDfjwuhfDasjdnJBDhfSDas"), nil
	})
	if err != nil {
		return "", err
	}
	if uid, ok := claim.Claims.(jwt.MapClaims)["username"].(string); ok {
		return uid, nil
	}
	return "", fmt.Errorf("fail parse")
}

func ExpireToken(token string) (bool, error) {
	//使token过期不能使用
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("iDfjwuhfDasjdnJBDhfSDas"), nil
	})
	if err != nil {
		return false, err
	}
	claim.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(time.Hour * -1).Unix()
	return true, nil
}

/*func ConvertFileToBytes(file *multipart.File) ([]byte, error) {
	// 获取文件名和文件类型
	fileHeader := file.Header.Get("Content-Disposition")
	fileNameStartIndex := strings.Index(fileHeader, `filename="`) + 10
	fileNameEndIndex := strings.Index(fileHeader, `"`)
	fileName := fileHeader[fileNameStartIndex:fileNameEndIndex]
	fileType := strings.ToLower(fileName[strings.Index(fileName, ".")+1:])

	// 读取文件内容
	fileBytes, err := ioutil.ReadAll(*file)
	if err != nil {
		return nil, err
	}

	// 根据文件类型进行解码
	var resultBytes []byte
	switch fileType {
	case "txt", "log", "md", "html":
		resultBytes = fileBytes
	case "pdf":
		pdfReader, err := model.NewPdfReader(bytes.NewReader(fileBytes))
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return nil, err
		}
		for i := 0; i < numPages; i++ {
			page, err := pdfReader.GetPage(i + 1)
			if err != nil {
				return nil, err
			}
			contentStreams, err := page.GetContentStreams()
			if err != nil {
				return nil, err
			}
			buf.WriteString(contentStreams)
			buf.WriteString("\n")
		}
		resultBytes = buf.Bytes()
	case "doc", "docx":
		doc, err := document.OpenFromReader(bytes.NewReader(fileBytes))
		if err != nil {
			return nil, err
		}
		resultBytes, err = doc.SaveToBuffer(document.Docx)
		if err != nil {
			return nil, err
		}
	case "xls", "xlsx":
		xlFile, err := xlsx.OpenBinary(fileBytes)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows {
				for _, cell := range row.Cells {
					text := cell.String()
					buf.WriteString(text)
					buf.WriteString("\t")
				}
				buf.WriteString("\n")
			}
		}
		resultBytes = buf.Bytes()
	default:
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}

	return resultBytes, nil
}
*/
