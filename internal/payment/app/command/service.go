package command

import (
	"context"
	"github.com/xh/gorder/internal/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
