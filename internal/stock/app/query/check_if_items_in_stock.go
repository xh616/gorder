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
	if err := h.checkStock(ctx, query.Items); err != nil {
		return nil, err
	}
	// TODO 扣库存
	return res, nil
}
func (h checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
	var ids []string
	for _, i := range query {
		ids = append(ids, i.ID)
	}
	records, err := h.stockRepo.GetStock(ctx, ids)
	if err != nil {
		return err
	}
	idQuantityMap := make(map[string]int32)
	for _, r := range records {
		idQuantityMap[r.ID] += r.Quantity
	}
	var (
		ok       = true
		failedOn []struct {
			ID   string
			Want int32
			Have int32
		}
	)
	for _, item := range query {
		if item.Quantity > idQuantityMap[item.ID] {
			ok = false
			failedOn = append(failedOn, struct {
				ID   string
				Want int32
				Have int32
			}{ID: item.ID, Want: item.Quantity, Have: idQuantityMap[item.ID]})
		}
	}
	if ok {
		return nil
	}
	return domain.ExceedStockError{FailedOn: failedOn}
}

//func getStubPriceID(id string) string {
//	priceID, ok := stub[id]
//	if !ok {
//		priceID = stub["1"]
//	}
//	return priceID
//}
