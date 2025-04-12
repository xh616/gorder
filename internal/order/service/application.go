package service

import (
	"context"
	"github.com/sirupsen/logrus"
	grpcClient "github.com/xh/gorder/internal/common/client"
	"github.com/xh/gorder/internal/common/metrics"
	"github.com/xh/gorder/internal/order/adapters"
	"github.com/xh/gorder/internal/order/adapters/grpc"
	"github.com/xh/gorder/internal/order/app"
	"github.com/xh/gorder/internal/order/app/command"
	"github.com/xh/gorder/internal/order/app/query"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockClient()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
