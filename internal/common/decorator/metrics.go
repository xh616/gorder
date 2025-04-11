package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type MetricsClient interface {
	Inc(key string, value int) //对指定键（key）的指标值增加 value，需实现
}

type queryMetricsDecorator[C, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	defer func() {
		end := time.Since(start)
		// 执行时间记录
		q.client.Inc(fmt.Sprintf("querys.%s.duration", actionName), int(end.Seconds()))
		// 成功次数和失败次数
		if err == nil {
			q.client.Inc(fmt.Sprintf("querys.%s.success", actionName), 1)
		} else {
			q.client.Inc(fmt.Sprintf("querys.%s.failure", actionName), 1)
		}
	}()
	return q.base.Handle(ctx, cmd)
}
