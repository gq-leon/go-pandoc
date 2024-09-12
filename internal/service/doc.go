// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IDoc interface {
		ToDocx(ctx context.Context, src string) (string, error)
		ToMarkdown(ctx context.Context, src string) (string, error)
	}
)

var (
	localDoc IDoc
)

func Doc() IDoc {
	if localDoc == nil {
		panic("implement not found for interface IDoc, forgot register?")
	}
	return localDoc
}

func RegisterDoc(i IDoc) {
	localDoc = i
}
