package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type DocReq struct {
	g.Meta `path:"/tools/doc2md" tags:"工具" mime:"multipart/form-data" method:"POST" summary:"文件上传转markdown"`
	Files  []*ghttp.UploadFile `json:"file[]" type:"file" dc:"选择上传文件"`
}

type DocRes struct {
	Urls []string `json:"urls" dc:"文件访问URL"`
}
