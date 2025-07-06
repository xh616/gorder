package logging

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// Init 日志输出格式配置
func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})
	// 要结构化日志就注释掉，不要的话就留着
	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
		//logger.SetFormatter(&prefixed.TextFormatter{
		//	ForceFormatting: true,
		//})
	}
}
