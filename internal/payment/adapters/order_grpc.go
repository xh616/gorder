package adapters

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xh/gorder/internal/common/genproto/orderpb"
	"google.golang.org/grpc/status"
)

type OrderGRPC struct {
	client orderpb.OrderServiceClient
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{client: client}
}

func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) (err error) {
	defer func() {
		if err != nil {
			logrus.Infof("payment_adapter||update_order,err=%v", err)
		}
	}()
	_, err = o.client.UpdateOrder(ctx, order)
	return status.Convert(err).Err()
}
