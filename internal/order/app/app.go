package app

import "github.com/xh/gorder/internal/order/app/query"

type Application struct {
	Commands
	Queries
}

type Commands struct {
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
