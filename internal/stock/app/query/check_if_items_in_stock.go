package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xh/gorder/internal/common/decorator"
	domain "github.com/xh/gorder/internal/stock/domain/stock"
	"github.com/xh/gorder/internal/stock/entity"
	"github.com/xh/gorder/internal/stock/infrastructure/integration"
)

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
	stripeAPI *integration.StripeAPI
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	stripeAPI *integration.StripeAPI,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	if stripeAPI == nil {
		panic("nil stripeAPI")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*entity.Item](
		checkIfItemsInStockHandler{
			stockRepo: stockRepo,
			stripeAPI: stripeAPI,
		},
		logger,
		metricClient,
	)
}

//var stub = map[string]string{
//	"1": "price_1RT1Y3POtxEUnZqFie3rULoy",
//}

func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
	var res []*entity.Item
	for _, item := range query.Items {
		//从stripe 获取
		priceID, err := h.stripeAPI.GetPriceByProductID(ctx, item.ID)
		if err != nil || priceID == "" {
			return nil, err
		}
		//priceID := getStubPriceID(item.ID)
		res = append(res, &entity.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}
	return res, nil
}

//func getStubPriceID(id string) string {
//	priceID, ok := stub[id]
//	if !ok {
//		priceID = stub["1"]
//	}
//	return priceID
//}
