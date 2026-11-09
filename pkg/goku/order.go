// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package goku

import (
	"fmt"
	"strings"
)

// OrderMethod represents the sorting direction (ASC or DESC).
type OrderMethod string

const (
	// OrderMethodASC represents ascending order.
	OrderMethodASC OrderMethod = "ASC"
	// OrderMethodDESC represents descending order.
	OrderMethodDESC OrderMethod = "DESC"
)

// Order defines a column and its sorting method.
type Order struct {
	Column string
	Method OrderMethod
}

// String returns the string representation of the OrderMethod.
func (m OrderMethod) String() string {
	return string(m)
}

// OrdersBuilder parses a comma-separated string into a slice of Order.
// Multiple columns are separated by commas ",".
func OrdersBuilder(str string) []Order {
	res := []Order{}
	if str == "" {
		return res
	}

	list := strings.SplitSeq(str, ",")
	for v := range list {
		if v == "" || v == "-" {
			continue
		}
		res = append(res, OrderBuilder(v))
	}

	return res
}

// OrderBuilder parses a single sorting string into an Order.
// Add a prefix "-" to indicate that the order method is DESC.
//
// Example:
//   - address -> {Column: "address", Method: ASC}
//   - -address -> {Column: "address", Method: DESC}
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

// OrdersQueryBuilder builds a SQL ORDER BY query string from a slice of Order
// using the provided column map to translate field names to actual database columns.
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
