package server

import (
	"yuanlimm-worker/api"
	"yuanlimm-worker/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Cors())

	// 路由
	v1 := r.Group("/api/")
	{
		v1.GET("ping", api.Ping)
		v1.GET("super_wishs", api.SuperWishsConfig)
		v1.POST("super_wishs", api.SuperWishs)
	}
	return r
}
