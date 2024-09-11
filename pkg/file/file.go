package file

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GetSize get the file size
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)

	return len(content), err
}

// GetExt get the file ext
func GetExt(fileName string) string {
	return filepath.Ext(fileName)
}

// CheckNotExist check if the file exists
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	if err == nil {
		return true
	}
	return os.IsNotExist(err)
}

// CheckPermission check if the file has permission
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// IsNotExistMkDir create a directory if it does not exist
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir create a directory
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Open a file according to a specific mode
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// MustOpen maximize trying to open the file
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

// ReplaceExt replace file ext
func ReplaceExt(filePath string, newExt string) string {
	if ext := GetExt(filePath); ext != "" {
		return strings.TrimSuffix(filePath, ext) + newExt
	}
	return filePath + newExt
}

// Hash 计算文件hash值
func Hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", path, err)
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes), nil
}

// ReplaceContent 文件内容替换
func ReplaceContent(filePath string, data map[string]string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}
	content := string(bytes)
	for src, dst := range data {
		log.Println(src, dst)
		re := regexp.MustCompile(`(\.?)` + src)
		content = re.ReplaceAllString(content, dst)
	}
	err = os.WriteFile(filePath, []byte(content), os.ModePerm)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	return nil
}

func GetContentType(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buffer := make([]byte, 512)

	_, err = f.Read(buffer)
	if err != nil {
		return ""
	}
	contentType := http.DetectContentType(buffer)
	return contentType
}
