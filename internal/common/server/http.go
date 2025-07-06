package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr")
	if addr == "" {
		panic("http address is empty")
	}
	RunHTTPServerOnAddr(addr, wrapper)
}

// RunHTTPServerOnAddr 路由注册和服务初始化
func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	setMiddlewares(apiRouter)
	wrapper(apiRouter)
	apiRouter.Group("/api")
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}

// setMiddlewares 设置中间件
func setMiddlewares(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("default_server"))
}
