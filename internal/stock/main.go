package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/xh/gorder/internal/common/config"
	"github.com/xh/gorder/internal/common/discovery"
	"github.com/xh/gorder/internal/common/genproto/stockpb"
	"github.com/xh/gorder/internal/common/logging"
	"github.com/xh/gorder/internal/common/server"
	"github.com/xh/gorder/internal/common/tracing"
	"github.com/xh/gorder/internal/stock/ports"
	"github.com/xh/gorder/internal/stock/service"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 初始化jaeger
	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application := service.NewApplication(ctx)

	// 注册consul
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer(application)
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		// 暂时不用
	default:
		panic("unexpected server type: " + serverType)
	}
}
