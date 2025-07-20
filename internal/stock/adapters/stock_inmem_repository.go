package adapters

import (
	"context"
	domain "github.com/xh/gorder/internal/stock/domain/stock"
	"github.com/xh/gorder/internal/stock/entity"
	"sync"
)

type MemoryStockRepository struct {
	lock  *sync.RWMutex
	store map[string]*entity.Item
}

var stub = map[string]*entity.Item{
	"item_id": {
		ID:       "foo_item",
		Name:     "stub item",
		Quantity: 10000,
		PriceID:  "stub_item_price_id",
	},
	"item1": {
		ID:       "item1",
		Name:     "stub item 1",
		Quantity: 10000,
		PriceID:  "stub_item1_price_id",
	},
	"item2": {
		ID:       "item2",
		Name:     "stub item 2",
		Quantity: 10000,
		PriceID:  "stub_item2_price_id",
	},
	"item3": {
		ID:       "item3",
		Name:     "stub item 3",
		Quantity: 10000,
		PriceID:  "stub_item3_price_id",
	},
}

func NewMemoryStockRepository() *MemoryStockRepository {
	return &MemoryStockRepository{
		lock:  new(sync.RWMutex),
		store: stub,
	}
}

func (m MemoryStockRepository) GetItems(_ context.Context, ids []string) ([]*entity.Item, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var (
		res     []*entity.Item
		missing []string
	)
	for _, id := range ids {
		if item, exist := m.store[id]; exist {
			res = append(res, item)
		} else {
			missing = append(missing, id)
		}
	}
	if len(res) == len(ids) {
		return res, nil
	}
	return res, domain.NotFoundError{Missing: missing}
}
