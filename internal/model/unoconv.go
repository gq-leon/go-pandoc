package model

import (
	"context"
)

type UnoconvCall struct {
	Ctx      context.Context
	Args     []string
	ReTry    int
	CallBack chan error
}

func NewUnoconvCall(ctx context.Context, args []string) *UnoconvCall {
	return &UnoconvCall{
		Ctx:      ctx,
		Args:     args,
		ReTry:    0,
		CallBack: make(chan error, 1),
	}
}
