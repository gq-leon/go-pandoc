// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"github.com/gq-leon/go-pandoc/internal/model"
)

type (
	IUnoconv interface {
		Add(data *model.UnoconvCall) error
	}
)

var (
	localUnoconv IUnoconv
)

func Unoconv() IUnoconv {
	if localUnoconv == nil {
		panic("implement not found for interface IUnoconv, forgot register?")
	}
	return localUnoconv
}

func RegisterUnoconv(i IUnoconv) {
	localUnoconv = i
}
