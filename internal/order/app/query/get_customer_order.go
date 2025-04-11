package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xh/gorder/internal/common/decorator"
	domain "github.com/xh/gorder/internal/order/domain/order"
)

type GetCustomerOrder struct {
	CustomerID string `json:"customer_id"`
	OrderID    string `json:"order_id"`
}

// GetCustomerOrderHandler 是QueryHandler别名
type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

// getCustomerOrderHandler实现了Handle，那么他就是一个QueryHandler了
type getCustomerOrderHandler struct {
	orderRepo domain.Repository
}

func NewGetCustomerOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
		getCustomerOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

// Handle 具体实现查询
func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
	if err != nil {
		return nil, err
	}
	return o, nil
}
