package server

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"net"

	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpctags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	grpclogrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
}

// RunGRPCServer find grpc config
func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
	addr := viper.Sub(serviceName).GetString("grpc-addr")
	if addr == "" {
		//TODO: Warning log
		addr = viper.GetString("fallback-grpc-addr")
	}
	RunGRPCServerOnAddr(addr, registerServer)
}

// RunGRPCServerOnAddr start grpc service
// 注入 grpc 的 生态 ecosystem
func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			grpctags.UnaryServerInterceptor(grpctags.WithFieldExtractor(grpctags.CodeGenRequestFieldExtractor)),
			grpclogrus.UnaryServerInterceptor(logrusEntry),
		),
		grpc.ChainStreamInterceptor(
			grpctags.StreamServerInterceptor(grpctags.WithFieldExtractor(grpctags.CodeGenRequestFieldExtractor)),
			grpclogrus.StreamServerInterceptor(logrusEntry),
		),
	)

	registerServer(grpcServer)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("Starting gRPC Server, Listening: %s", addr)
	if err := grpcServer.Serve(listen); err != nil {
		logrus.Panic(err)
	}
}
