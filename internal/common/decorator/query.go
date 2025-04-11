package decorator

import (
	"context"
	"github.com/sirupsen/logrus"
)

// QueryHandler defines a generic type that receives a Query Q,
// and returns a result R
type QueryHandler[Q, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
	// 日志增强，queryLoggingDecorator实现了QueryHandler接口
	return queryLoggingDecorator[H, R]{
		logger: logger,
		// metrics增强
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
