package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xh/gorder/internal/common/metrics"
	"github.com/xh/gorder/internal/order/adapters"
	"github.com/xh/gorder/internal/order/app"
	"github.com/xh/gorder/internal/order/app/query"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
