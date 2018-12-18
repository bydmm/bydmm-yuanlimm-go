package api

import (
	"net/http"
	"yuanlimm-worker/serializer"
	"yuanlimm-worker/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
)

// Ping 状态检查页面
func Ping(c *gin.Context) {
	c.String(200, "Pong")
}

// SuperWishsConfig 许愿配置
func SuperWishsConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hard":      24,
		"unix_time": now.BeginningOfMinute().Unix(),
	})
}

// SuperWishs 许愿
func SuperWishs(c *gin.Context) {
	service := service.SuperWishsService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Wish()
		c.JSON(200, res)
	} else {
		c.JSON(200, serializer.WishResponse{
			Success: false,
			Hard:    24,
			Msg:     "应援失败，再接再厉！",
		})
	}
}
