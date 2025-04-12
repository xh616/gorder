package app

import (
	"github.com/xh/gorder/internal/stock/app/query"
)

type Application struct {
	Commands
	Queries
}

type Commands struct {
}

type Queries struct {
	CheckIfItemsInStock query.CheckIfItemsInStockHandler
	GetItems            query.GetItemsHandler
}
