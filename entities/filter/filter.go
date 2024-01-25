// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package filter

import "strings"

type OrderMethod string

const (
	OrderMethodASC  OrderMethod = "ASC"
	OrderMethodDESC OrderMethod = "DESC"
)

type Filter struct {
	Search string
	Page   int
	Limit  int
	Orders []Order
}

type Order struct {
	Column string
	Method OrderMethod
}

func (filter *Filter) OrdersBuilder(str string) {
	if str == "" {
		return
	}
	list := strings.Split(str, ",")
	for _, v := range list {
		if v == "" || v == "-" {
			continue
		}
		filter.Orders = append(filter.Orders, OrderBuilder(v))
	}
}

func OrderBuilder(str string) Order {
	if str[0:1] == "-" {
		return Order{
			Column: str[1:],
			Method: OrderMethodDESC,
		}
	} else {
		return Order{
			Column: str,
			Method: OrderMethodASC,
		}
	}
}
