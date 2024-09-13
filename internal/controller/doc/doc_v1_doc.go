package doc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gq-leon/go-pandoc/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "github.com/gq-leon/go-pandoc/api/doc/v1"
)

func (c *ControllerV1) Doc(ctx context.Context, req *v1.DocReq) (res *v1.DocRes, err error) {
	if req.Files == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "请选择需要上传的文件")
	}

	var uploadFolder = "./uploads/" + time.Now().Format("2006/01/02")
	var filenames []string
	for _, file := range req.Files {
		filename, err := file.Save(uploadFolder, true)
		if err != nil {
			g.Log().Errorf(ctx, "文件上传失败:%s", err)
			return nil, gerror.NewCode(gcode.CodeMissingParameter, "文件上传失败")
		}
		g.Log().Infof(ctx, "文件上传成功 %v", filename)
		filenames = append(filenames, filename)
	}

	var (
		urls    []string
		wg      sync.WaitGroup
		mu      sync.Mutex
		errChan = make(chan error, 1)
	)

	for _, filename := range filenames {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			dst, err := service.Doc().ToMarkdown(gctx.New(), fmt.Sprintf("%s/%s", uploadFolder, filename))
			if err != nil {
				errChan <- gerror.NewCode(gcode.CodeInternalError, err.Error())
				return
			}
			mu.TryLock()
			urls = append(urls, dst)
			mu.Unlock()
		}(filename)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if finalErr, ok := <-errChan; ok {
		return nil, finalErr
	}

	return &v1.DocRes{Urls: urls}, nil
}
