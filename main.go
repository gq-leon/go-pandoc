package main

import (
	_ "github.com/gq-leon/go-pandoc/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gq-leon/go-pandoc/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
