package e

var MsgFlags = map[int]string{
	Success:         "ok",
	Error:           "fail",
	ErrorFileUpload: "文件上传失败",
	ErrorFileCover:  "文件格式转换失败",
	InvalidParams:   "请求参数错误",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Error]
}
