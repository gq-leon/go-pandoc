package pandoc_service

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	fileUtils "pandoc/pkg/file"
	"pandoc/pkg/logging"
	"pandoc/pkg/utils"
)

type Options struct {
	ExtractMedia string
	Dst          string // 输出文件地址
	Src          string // 源文件地址
	DownloadUrl  string // 文件下载地址

	mu   sync.Mutex
	args []string
}

func (s *Options) Run() error {
	if err := s.cover(); err != nil {
		return err
	}
	return s.resourceProcess()
}

func (s *Options) cover() error {
	s.spliceArgs()
	var (
		stdout  bytes.Buffer
		stderr  bytes.Buffer
		command = exec.Command("pandoc", s.args...)
	)
	logging.Debug(fmt.Sprintf("run command: pandoc " + strings.Join(s.args, " ")))

	// 打印输出内容和错误信息
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if s := stderr.String(); s != "" {
		logging.Error(fmt.Sprintf("pandoc err: %s", stderr.String()))
	}
	return err
}

func (s *Options) spliceArgs() {
	if s.ExtractMedia != "" {
		s.args = append(s.args, "--extract-media", s.ExtractMedia)
	}
	s.args = append(s.args, "--output", s.Dst, s.Src)
}

func (s *Options) resourceProcess() error {
	if !fileUtils.CheckNotExist(s.ExtractMedia) {
		return nil
	}
	logging.Debug("start resource process...")
	var (
		files = make(map[string]string) // map[local file path][cloud url]
		wg    sync.WaitGroup
	)
	if err := filepath.Walk(s.ExtractMedia+"/media", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			if upload, err := utils.OssUpload("pandoc", filePath); err == nil {
				s.mu.Lock()
				files[filePath] = fmt.Sprintf("%s/%s", utils.OssUrl, upload)
				s.mu.Unlock()
			} else {
				logging.Error(err)
			}
		}(path)
		return nil
	}); err != nil {
		return err
	}
	wg.Wait()

	go func() {
		if err := os.RemoveAll(s.ExtractMedia); err != nil {
			logging.Error(fmt.Sprintf("删除资源失败: %s", err))
		}
	}()

	logging.Debug("资源地址替换...")
	if err := fileUtils.ReplaceContent(s.Dst, files); err != nil {
		logging.Error(fmt.Sprintf("资源地址替换: %s", err))
	}

	// 文件上传
	logging.Debug("上传生成文件...")
	if ossUpload, err := utils.OssUpload("pandoc/markdown", s.Dst); err != nil {
		logging.Error(fmt.Sprintf("生成文件上传失败: %s", err))
	} else {
		s.DownloadUrl = fmt.Sprintf("%s/%s", utils.OssUrl, ossUpload)
	}

	return nil
}
