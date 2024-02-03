// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package goku

import (
	"fmt"
	"strings"
)

type OrderMethod string

const (
	OrderMethodASC  OrderMethod = "ASC"
	OrderMethodDESC OrderMethod = "DESC"
)

type Order struct {
	Column string
	Method OrderMethod
}

func (m OrderMethod) String() string {
	return string(m)
}

// build []Order by string. multiple columns separated by comma ","
func OrdersBuilder(str string) []Order {
	res := []Order{}
	if str == "" {
		return res
	}
	list := strings.Split(str, ",")
	for _, v := range list {
		if v == "" || v == "-" {
			continue
		}
		res = append(res, OrderBuilder(v))
	}
	return res
}

// build Order by string. add prefix "-" to indicate that order method is DESC.
//
// example:
//
// address -> {Column: address, Method: ASC}
//
// -address -> {Column: address, Method: DESC}
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

func OrdersQueryBuilder(orders []Order, columnMap map[string]string) (queryString string, err error) {
	ordersLen := len(orders)
	if ordersLen == 0 || (ordersLen == 0 && len(columnMap) == 0) {
		return "", nil
	}
	columns := make([]string, 0, ordersLen)
	for _, v := range orders {
		if column := columnMap[v.Column]; column != "" {
			columns = append(columns, column+" "+v.Method.String())
		} else {
			return "", fmt.Errorf("column %s not found", v.Column)
		}
	}
	return strings.Join(columns, ","), nil
}
