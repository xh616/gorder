package decorator

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

// 在QueryHandler上封装日志
type queryLoggingDecorator[C, R any] struct {
	logger *logrus.Entry
	base   QueryHandler[C, R]
}

// Handle 实现QueryHandler接口
func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := q.logger.WithFields(logrus.Fields{
		"query":      generateActionName(cmd),
		"query_body": fmt.Sprintf("%#v", cmd), //%#v输出带类型信息，适合调试
	})
	logger.Debug("Executing query")
	defer func() {
		if err == nil {
			logger.Info("Query execute successfully")
		} else {
			logger.Error("Failed to execute query", err)
		}
	}()
	return q.base.Handle(ctx, cmd)
}

func generateActionName(cmd any) string {
	// e.g. query.XXXHandle ==> XXXHandle
	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
}
