package domain

import (
	"context"
	"github.com/xh/gorder/internal/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}
