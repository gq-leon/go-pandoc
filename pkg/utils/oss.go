package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"pandoc/pkg/file"
	"pandoc/pkg/logging"
)

var (
	bucket *oss.Bucket
	OssUrl string
)

func init() {
	log.Println("start init oss...")
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		panic(err)
	}
	endpoint := os.Getenv("OSS_ENDPOINT")
	if endpoint == "" {
		endpoint = "oss-cn-chengdu.aliyuncs.com"
	}
	client, err := oss.New(endpoint, "", "", oss.SetCredentialsProvider(&provider))
	if err != nil {
		panic(err)
	}
	bucket, err = client.Bucket("agi-apq")
	if err != nil {
		panic(err)
	}

	if OssUrl = os.Getenv("OSS_URL"); OssUrl == "" {
		OssUrl = "https://agi-apq.oss-cn-chengdu.aliyuncs.com"
	}
	log.Println("start init oss success.")
}

// OssUpload 上传返回文件地址
func OssUpload(path, localPath string) (string, error) {
	hash, _ := file.Hash(localPath)
	objectKey := fmt.Sprintf("%s/%s%s", path, hash, file.GetExt(localPath))
	if err := bucket.PutObjectFromFile(objectKey, localPath); err != nil {
		return "", err
	}
	if err := bucket.SetObjectMeta(objectKey, oss.ContentType(file.GetContentType(localPath))); err != nil {
		logging.Error("set file content type err", err)
	}
	return objectKey, nil
}
