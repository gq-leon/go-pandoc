package oss

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gq-leon/go-pandoc/internal/service"
	"net/http"
)

var (
	bucket *oss.Bucket
	ossUrl string
)

type sOss struct {
}

func init() {
	ctx := gctx.New()
	g.Log().Info(ctx, "initializing oss...")
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		panic(err)
	}

	endpoint, _ := g.Cfg().GetWithEnv(ctx, "oss.endpoint", "oss-cn-chengdu.aliyuncs.com")
	bucketName, _ := g.Cfg().GetWithEnv(ctx, "oss.bucketName", "agi-apq")
	url, _ := g.Cfg().GetWithEnv(ctx, "oss.url", "https://agi-apq.oss-cn-chengdu.aliyuncs.com")
	ossUrl = url.String()

	client, err := oss.New(endpoint.String(), "", "", oss.SetCredentialsProvider(&provider))
	if err != nil {
		panic(err)
	}

	if bucket, err = client.Bucket(bucketName.String()); err != nil {
		panic(err)
	}

	service.RegisterOss(New())
}

func New() service.IOss {
	return &sOss{}
}

func (s *sOss) Upload(ctx context.Context, path, localPath string) (string, error) {
	hash, _ := gmd5.EncryptFile(localPath)
	objectKey := fmt.Sprintf("%s/%s%s", path, hash, gfile.Ext(localPath))
	if err := bucket.PutObjectFromFile(objectKey, localPath); err != nil {
		return "", err
	}

	bytes := gfile.GetBytesByTwoOffsetsByPath(localPath, 0, 512)
	if err := bucket.SetObjectMeta(objectKey, oss.ContentType(http.DetectContentType(bytes))); err != nil {
		g.Log().Errorf(ctx, "set file content type err:%s", err)
	}
	return fmt.Sprintf("%s/%s", ossUrl, objectKey), nil
}
