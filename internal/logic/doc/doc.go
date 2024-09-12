package doc

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gq-leon/go-pandoc/internal/service"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type sDoc struct{}

func init() {
	service.RegisterDoc(New())
}

func New() service.IDoc {
	return &sDoc{}
}

func (s *sDoc) ToDocx(ctx context.Context, src string) (string, error) {
	//	参考 https://en.wikipedia.org/wiki/List_of_file_signatures
	headers := gfile.GetBytesByTwoOffsetsByPath(src, 0, 4)
	if headers[0] == 0xD0 && headers[1] == 0xCF && headers[2] == 0x11 && headers[3] == 0xE0 {
		var (
			dst  = fmt.Sprintf("%s/%s.docx", gfile.Dir(src), gfile.Name(src))
			args = []string{"-f", "docx", "-o", dst, src}
		)
		g.Log().Infof(ctx, "exec command: unoconv %s", strings.Join(args, " "))
		command := exec.Command("unoconv", args...)
		if _, err := command.Output(); err != nil {
			return "", err
		}
		return dst, nil
	}

	return src, nil
}

func (s *sDoc) ToMarkdown(ctx context.Context, src string) (string, error) {
	var (
		dst          = fmt.Sprintf("%s/%s.md", gfile.Dir(src), gfile.Name(src))
		extractMedia = fmt.Sprintf("%s/%s", gfile.Dir(src), guid.S())
		err          error
	)
	if src, err = s.ToDocx(ctx, src); err != nil {
		g.Log().Errorf(ctx, "ToDocx err: %s", err)
		return "", errors.New("文件转docx失败")
	}

	args := []string{"--extract-media", extractMedia, "--output", dst, src}
	g.Log().Infof(ctx, "exec command: pandoc %s", strings.Join(args, " "))
	command := exec.Command("pandoc", args...)
	if _, err := command.Output(); err != nil {
		return "", errors.New("文件类型转换失败")
	}

	if resource := mediaResource(ctx, extractMedia); len(resource) > 0 {
		g.Log().Info(ctx, "媒体资源替换...")
		if err := replaceMedia(dst, resource); err != nil {
			g.Log().Errorf(ctx, "媒体资源替换失败: %s", err)
		}
	}

	upload, err := service.Oss().Upload(ctx, "pandoc", dst)
	if err != nil {
		return "", err
	}
	return upload, nil
}

func mediaResource(ctx context.Context, path string) map[string]string {
	files := make(map[string]string)
	mediaPath := path + "/media"
	if gfile.Exists(mediaPath) {
		list, _ := gfile.ScanDir(mediaPath, "*", false)
		for _, v := range list {
			if upload, err := service.Oss().Upload(ctx, "pandoc", v); err != nil {
				g.Log().Errorf(ctx, "上传 oss 失败: %s", err)
			} else {
				abs, _ := filepath.Abs(".")
				rel, _ := filepath.Rel(abs, v)
				files[rel] = upload
			}
		}
		_ = gfile.Remove(path)
	}
	return files
}

func replaceMedia(filePath string, data map[string]string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}
	content := string(bytes)
	for src, dst := range data {
		re := regexp.MustCompile(`(\.?)` + src)
		content = re.ReplaceAllString(content, dst)
	}
	err = os.WriteFile(filePath, []byte(content), os.ModePerm)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	return nil
}
