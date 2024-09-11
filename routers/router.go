package routers

import (
	"github.com/gin-gonic/gin"

	v1 "pandoc/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/doc2md", v1.Doc2Md)
	}
	return r
}
