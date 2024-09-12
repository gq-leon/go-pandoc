// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IOss interface {
		Upload(ctx context.Context, path string, localPath string) (string, error)
	}
)

var (
	localOss IOss
)

func Oss() IOss {
	if localOss == nil {
		panic("implement not found for interface IOss, forgot register?")
	}
	return localOss
}

func RegisterOss(i IOss) {
	localOss = i
}
