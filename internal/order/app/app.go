package app

import (
	"github.com/xh/gorder/internal/order/app/command"
	"github.com/xh/gorder/internal/order/app/query"
)

type Application struct {
	Commands
	Queries
}

type Commands struct {
	CreateOrder command.CreateOrderHandler
	UpdateOrder command.UpdateOrderHandler
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
