package client

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xh/gorder/internal/common/discovery"
	"github.com/xh/gorder/internal/common/genproto/orderpb"
	"github.com/xh/gorder/internal/common/genproto/stockpb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"time"
)

func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
	if !WaitForOrderGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
		return nil, nil, errors.New("order grpc client not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, func() error { return nil }, err
	}
	if grpcAddr == "" {
		logrus.Warn("empty grpc addr for order grpc")
	}
	opts := grpcDialOpts(grpcAddr)

	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	if !WaitForStockGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
		return nil, nil, errors.New("stock grpc client not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, func() error { return nil }, err
	}
	if grpcAddr == "" {
		logrus.Warn("empty grpc addr for stock grpc")
	}
	opts := grpcDialOpts(grpcAddr)
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(_ string) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
}

func WaitForOrderGRPCClient(timeout time.Duration) bool {
	logrus.Infof("waiting for order grpc client, timeout: %v seconds", timeout.Seconds())
	return waitFor(viper.GetString("order.grpc-addr"), timeout)
}

func WaitForStockGRPCClient(timeout time.Duration) bool {
	logrus.Infof("waiting for stock grpc client, timeout: %v seconds", timeout.Seconds())
	return waitFor(viper.GetString("stock.grpc-addr"), timeout)
}

// waitFor 等待直到端口可用
func waitFor(addr string, timeout time.Duration) bool {
	portAvailable := make(chan struct{})
	timeoutCh := time.After(timeout)

	go func() {
		for {
			select {
			case <-timeoutCh:
				return
			default:
				// continue
			}
			_, err := net.Dial("tcp", addr) //ping
			if err == nil {                 // port is available
				close(portAvailable) //关闭的时候自动返回0
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {
	case <-portAvailable:
		return true
	case <-timeoutCh:
		return false
	}
}
