package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xh/gorder/internal/common/metrics"
	"github.com/xh/gorder/internal/stock/adapters"
	"github.com/xh/gorder/internal/stock/app"
	"github.com/xh/gorder/internal/stock/app/query"
	"github.com/xh/gorder/internal/stock/infrastructure/integration"
)

func NewApplication(ctx context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	stripeAPI := integration.NewStripeAPI()
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}
