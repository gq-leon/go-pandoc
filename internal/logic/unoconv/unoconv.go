package unoconv

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gq-leon/go-pandoc/internal/model"
	"github.com/gq-leon/go-pandoc/internal/service"
)

type sUnoconv struct {
	stacks chan *model.UnoconvCall
}

func init() {
	service.RegisterUnoconv(New())
}

func New() service.IUnoconv {
	service := &sUnoconv{
		stacks: make(chan *model.UnoconvCall),
	}
	if err := service.run(); err != nil {
		panic(err)
	}
	return service
}

func (s *sUnoconv) run() error {
	g.Log().Info(gctx.New(), "start unoconv stacks...")

	go func() {
		for {
			select {
			case data := <-s.stacks:
				// check retry num
				if data.ReTry > 3 {
					data.CallBack <- errors.New("stacks unoconv 执行已尝试%d次, 依旧无法将文件成功转换，请检查文件数据")
					g.Log().Warningf(data.Ctx, "执行unoconv参数 %s", data.Args)
					break
				}

				g.Log().Infof(data.Ctx, "stacks exec command: unoconv %s try %d", strings.Join(data.Args, " "), data.ReTry)
				command := exec.Command("unoconv", data.Args...)
				if _, err := command.Output(); err != nil {
					data.ReTry += 1
					s.Add(data)
					g.Log().Warningf(data.Ctx, "stacks unoconv err:%s", err)
				} else {
					data.CallBack <- nil
				}
			}
		}
	}()
	return nil
}

func (s *sUnoconv) Add(data *model.UnoconvCall) error {
	s.stacks <- data
	return nil
}
