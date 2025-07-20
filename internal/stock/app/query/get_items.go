package query

import (
	"context"
	"github.com/xh/gorder/internal/stock/entity"

	"github.com/sirupsen/logrus"
	"github.com/xh/gorder/internal/common/decorator"
	domain "github.com/xh/gorder/internal/stock/domain/stock"
)

type GetItems struct {
	ItemIDs []string
}

type GetItemsHandler decorator.QueryHandler[GetItems, []*entity.Item]

type getItemsHandler struct {
	stockRepo domain.Repository
}

func NewGetItemsHandler(
	stockRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) GetItemsHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	return decorator.ApplyQueryDecorators[GetItems, []*entity.Item](
		getItemsHandler{stockRepo: stockRepo},
		logger,
		metricClient,
	)
}

func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*entity.Item, error) {
	items, err := g.stockRepo.GetItems(ctx, query.ItemIDs)
	if err != nil {
		return nil, err
	}
	return items, nil
}
