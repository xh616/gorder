package query

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xh/gorder/internal/common/decorator"
	"github.com/xh/gorder/internal/common/handler/redis"
	domain "github.com/xh/gorder/internal/stock/domain/stock"
	"github.com/xh/gorder/internal/stock/entity"
	"github.com/xh/gorder/internal/stock/infrastructure/integration"
	"strings"
	"time"
)

const (
	redisLockPrefix = "check_stock_"
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
	// 用redis加分布式锁
	if err := lock(ctx, getLockKey(query)); err != nil {
		return nil, errors.Wrapf(err, "redis lock error: key=%s", getLockKey(query))
	}
	defer func() {
		if err := unlock(ctx, getLockKey(query)); err != nil {
			logrus.Warnf("redis unlock fail, err=%v", err)
		}
	}()

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
	return res, nil
}

func getLockKey(query CheckIfItemsInStock) string {
	var ids []string
	for _, i := range query.Items {
		ids = append(ids, i.ID)
	}
	return redisLockPrefix + strings.Join(ids, "_")
}

func lock(ctx context.Context, key string) error {
	return redis.SetNX(ctx, redis.LocalClient(), key, "1", 5*time.Minute)
}

func unlock(ctx context.Context, key string) error {
	return redis.Del(ctx, redis.LocalClient(), key)
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
		return h.stockRepo.UpdateStock(ctx, query, func(
			ctx context.Context,
			existing []*entity.ItemWithQuantity,
			query []*entity.ItemWithQuantity,
		) ([]*entity.ItemWithQuantity, error) {
			var newItems []*entity.ItemWithQuantity
			for _, e := range existing {
				for _, q := range query {
					if e.ID == q.ID {
						newItems = append(newItems, &entity.ItemWithQuantity{
							ID:       e.ID,
							Quantity: e.Quantity - q.Quantity,
						})
					}
				}
			}
			return newItems, nil
		})
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
