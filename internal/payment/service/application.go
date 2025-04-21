package service

import (
	"context"
	"github.com/sirupsen/logrus"
	grpcClient "github.com/xh/gorder/internal/common/client"
	"github.com/xh/gorder/internal/common/metrics"
	"github.com/xh/gorder/internal/payment/adapters"
	"github.com/xh/gorder/internal/payment/app"
	"github.com/xh/gorder/internal/payment/app/command"
	"github.com/xh/gorder/internal/payment/domain"
	"github.com/xh/gorder/internal/payment/infrastructure/processor"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	memoryProcessor := processor.NewInmemProcessor()
	return newApplication(ctx, orderGRPC, memoryProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}
