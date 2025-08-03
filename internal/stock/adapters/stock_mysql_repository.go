package adapters

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xh/gorder/internal/stock/entity"
	"github.com/xh/gorder/internal/stock/infrastructure/persistent"
	"github.com/xh/gorder/internal/stock/infrastructure/persistent/builder"
	"gorm.io/gorm"
)

type MySQLStockRepository struct {
	db *persistent.MySQL
}

func NewMySQLStockRepository(db *persistent.MySQL) *MySQLStockRepository {
	return &MySQLStockRepository{db: db}
}

func (m MySQLStockRepository) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
	data, err := m.db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(ids...))
	if err != nil {
		return nil, errors.Wrap(err, "BatchGetStockByID error")
	}
	var result []*entity.ItemWithQuantity
	for _, d := range data {
		result = append(result, &entity.ItemWithQuantity{
			ID:       d.ProductID,
			Quantity: d.Quantity,
		})
	}
	return result, nil
}

func (m MySQLStockRepository) UpdateStock(
	ctx context.Context,
	data []*entity.ItemWithQuantity,
	updateFn func(
		ctx context.Context,
		existing []*entity.ItemWithQuantity,
		query []*entity.ItemWithQuantity,
	) ([]*entity.ItemWithQuantity, error),
) error {
	return m.db.StartTransaction(func(tx *gorm.DB) (err error) {
		defer func() {
			if err != nil {
				logrus.Warnf("update stock transaction err=%v", err)
			}
		}()
		// 悲观锁
		err = m.updatePessimistic(ctx, tx, data, updateFn)
		// 乐观锁
		//err = m.updateOptimistic(ctx, tx, data, updateFn)
		return err
	})
}

// 乐观锁update，先读version，后update，如果version不一致，则更新失败
func (m MySQLStockRepository) updateOptimistic(
	ctx context.Context,
	tx *gorm.DB,
	data []*entity.ItemWithQuantity,
	updateFn func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity,
	) ([]*entity.ItemWithQuantity, error)) error {
	for _, queryData := range data {
		var newestRecord *persistent.StockModel
		newestRecord, err := m.db.GetStockByID(ctx, builder.NewStock().ProductIDs(queryData.ID))
		if err != nil {
			return err
		}
		if err = m.db.Update(
			ctx,
			tx,
			builder.NewStock().ProductIDs(queryData.ID).Versions(newestRecord.Version).QuantityGT(queryData.Quantity),
			map[string]any{
				"quantity": gorm.Expr("quantity - ?", queryData.Quantity), //gorm的扣库存写法
				"version":  newestRecord.Version + 1,                      //成功update把version+1
			}); err != nil {
			return err
		}
	}

	return nil
}

func (m MySQLStockRepository) unmarshalFromDatabase(dest []persistent.StockModel) []*entity.ItemWithQuantity {
	var result []*entity.ItemWithQuantity
	for _, i := range dest {
		result = append(result, &entity.ItemWithQuantity{
			ID:       i.ProductID,
			Quantity: i.Quantity,
		})
	}
	return result
}

// 悲观锁update（for update 互斥锁）
func (m MySQLStockRepository) updatePessimistic(
	ctx context.Context,
	tx *gorm.DB,
	data []*entity.ItemWithQuantity,
	updateFn func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity,
	) ([]*entity.ItemWithQuantity, error)) error {
	var dest []persistent.StockModel
	dest, err := m.db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(getIDFromEntities(data)...).ForUpdate())
	if err != nil {
		return errors.Wrap(err, "failed to find data")
	}

	existing := m.unmarshalFromDatabase(dest)
	updated, err := updateFn(ctx, existing, data)
	if err != nil {
		return err
	}

	for _, upd := range updated {
		for _, query := range data {
			if upd.ID != query.ID {
				continue
			}
			if err = m.db.Update(ctx, tx, builder.NewStock().ProductIDs(upd.ID).QuantityGT(query.Quantity),
				map[string]any{"quantity": gorm.Expr("quantity - ?", query.Quantity)}); err != nil {
				return errors.Wrapf(err, "unable to update %s", upd.ID)
			}
		}
	}
	return nil
}

func getIDFromEntities(items []*entity.ItemWithQuantity) []string {
	var ids []string
	for _, i := range items {
		ids = append(ids, i.ID)
	}
	return ids
}
