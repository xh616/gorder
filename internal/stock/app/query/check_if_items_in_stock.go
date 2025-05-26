package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xh/gorder/internal/common/decorator"
	"github.com/xh/gorder/internal/common/genproto/orderpb"
	domain "github.com/xh/gorder/internal/stock/domain/stock"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
		checkIfItemsInStockHandler{stockRepo: stockRepo},
		logger,
		metricClient,
	)
}

// TODO: 后面删掉
var stub = map[string]string{
	"1": "price_1RT1Y3POtxEUnZqFie3rULoy",
}

func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var res []*orderpb.Item
	for _, item := range query.Items {
		//TODO: 改成从数据库 or stripe 获取
		priceId, ok := stub[item.ID]
		if !ok {
			priceId = stub["1"]
		}
		res = append(res, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  priceId,
		})
	}
	return res, nil
}
