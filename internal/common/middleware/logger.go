package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func StructuredLog(l *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// c.Next()之前代表在请求处理之前的处理
		c.Next()
		// c.Next()之后代表在请求处理完之后的响应
		elapsed := time.Since(t)
		l.WithFields(logrus.Fields{
			"request_url":     c.Request.RequestURI,
			"client_ip":       c.ClientIP(),
			"time_elapsed_ms": elapsed.Milliseconds(),
		}).Info("request out")
	}
}
