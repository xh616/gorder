package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr")
	//if addr == "" {
	//	// TODO: Warning log
	//}
	RunHTTPServerOnAddr(addr, wrapper)
}

// RunHTTPServerOnAddr 路由注册和服务初始化
func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	wrapper(apiRouter)
	apiRouter.Group("/api")
	apiRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong!!!")
	})
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}
