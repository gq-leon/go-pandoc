package v1

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"pandoc/pkg/app"
	"pandoc/pkg/e"
	fileUtils "pandoc/pkg/file"
	"pandoc/pkg/logging"
	"pandoc/service/pandoc_service"
)

var uploadPath = "./uploads"

// Doc2Md 支持多doc文件转markdown
func Doc2Md(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	form, err := ctx.MultipartForm()
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.Error, nil)
		return
	}

	var (
		currentTime = time.Now()
		folderPath  = filepath.Join(uploadPath, fmt.Sprintf("%d", currentTime.Year()), fmt.Sprintf("%02d", currentTime.Month()), fmt.Sprintf("%02d", currentTime.Day()))
	)

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		logging.Warn("创建文件夹失败:", err)
		appG.Response(http.StatusInternalServerError, e.Error, nil)
		return
	}

	var (
		temp     = time.Now().UnixNano()
		files    = form.File["file[]"]
		download []string
		mu       sync.Mutex
		wg       sync.WaitGroup
	)

	for i, file := range files {
		wg.Add(1)
		go func(index int, fileHeader *multipart.FileHeader) {
			defer wg.Done()
			var (
				src       = fmt.Sprintf("%s/%d-%s", folderPath, temp, fileHeader.Filename)
				mediaPath = fmt.Sprintf("%s/%d%d", folderPath, temp, index)
				dst       = fileUtils.ReplaceExt(src, ".md")
			)
			if err = ctx.SaveUploadedFile(fileHeader, src); err != nil {
				logging.Warn(err)
				return
			}
			logging.Debug(fmt.Sprintf("上传文件: %s, 文件路径: %s", fileHeader.Filename, src))

			service := &pandoc_service.Options{
				ExtractMedia: mediaPath,
				Dst:          dst,
				Src:          src,
			}
			if err := service.Run(); err != nil {
				logging.Warn(err)
				return
			}

			mu.Lock()
			download = append(download, service.DownloadUrl)
			mu.Unlock()
		}(i, file)
	}

	wg.Wait()
	appG.Response(http.StatusOK, e.Success, gin.H{"urls": download})
}
