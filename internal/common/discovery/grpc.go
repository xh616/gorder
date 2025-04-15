package discovery

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xh/gorder/internal/common/discovery/consul"
	"time"
)

func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
	// 创建 Consul 注册客户端
	registry, err := consul.New(viper.GetString("consul.addr")) //默认127.0.0.1:8500
	if err != nil {
		return func() error { return nil }, err
	}
	instanceID := GenerateInstanceID(serviceName)
	// 获取服务gRPC地址
	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")
	//注册服务到Consul
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		return func() error { return nil }, err
	}
	//启动健康检查协程
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logrus.Panicf("no heartbeat from %s to registry, err=%v", serviceName, err)
			}
			time.Sleep(time.Second)
		}
	}()
	//记录注册成功日志
	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        grpcAddr,
	}).Info("registered to consul")
	// 返回注销函数, 调用时将从 Consul 注销该服务，通过返回注销闭包，确保资源可以正确释放
	// 闭包的特性：捕获外部变量：可以访问定义它的函数内的变量
	//           保持变量状态：即使RegisterToConsul函数已经返回，这些被捕获的变量依然存在，当主调方调用返回的闭包函数时，它仍然可以访问这些变量
	//			 独立实例：每次外部函数调用都会创建新的闭包实例
	return func() error {
		return registry.Deregister(ctx, instanceID, serviceName)
	}, nil
}
