package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xh/gorder/internal/common/broker"
	grpcClient "github.com/xh/gorder/internal/common/client"
	"github.com/xh/gorder/internal/common/logging"
	"github.com/xh/gorder/internal/common/tracing"
	"github.com/xh/gorder/internal/kitchen/adapters"
	"github.com/xh/gorder/internal/kitchen/infrastructure/consumer"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("kitchen.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	orderClient, closeFunc, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	defer closeFunc()

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	orderGRPC := adapters.NewOrderGRPC(orderClient)
	go consumer.NewConsumer(orderGRPC).Listen(ch)

	// 接收信号，退出当前服务
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigs
		logrus.Infof("receive signal, exiting...")
		os.Exit(0)
	}()
	logrus.Println("to exit, press ctrl+c")
	select {}
}
