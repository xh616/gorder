package order

import (
	"context"
	"fmt"
)

type Repository interface {
	Create(context.Context, *Order) (*Order, error)
	Get(ctx context.Context, id, customerID string) (*Order, error)
	Update(
		ctx context.Context,
		o *Order,
		updateFun func(context.Context, *Order) (*Order, error),
	) error
}

type NotFoundError struct {
	OrderId string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("order not found: %s", e.OrderId)
}
