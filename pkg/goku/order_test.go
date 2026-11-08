// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package goku

import (
	"reflect"
	"testing"
)

const (
	name      = "name"
	createdAt = "created_at"
)

func TestOrdersBuilder(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want []Order
	}{
		{
			name: "single ASC and DESC columns",
			str:  "name,-created_at",
			want: []Order{
				{Column: name, Method: OrderMethodASC},
				{Column: createdAt, Method: OrderMethodDESC},
			},
		},
		{
			name: "empty string input",
			str:  "",
			want: []Order{},
		},
		{
			name: "ignore empty or dash items",
			str:  "name,,-,age",
			want: []Order{
				{Column: name, Method: OrderMethodASC},
				{Column: "age", Method: OrderMethodASC},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OrdersBuilder(tt.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrdersBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderBuilder(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want Order
	}{
		{
			name: "ascending order",
			str:  "email",
			want: Order{Column: "email", Method: OrderMethodASC},
		},
		{
			name: "descending order with minus prefix",
			str:  "-created_at",
			want: Order{Column: createdAt, Method: OrderMethodDESC},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OrderBuilder(tt.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdersQueryBuilder(t *testing.T) {
	columnMap := map[string]string{
		"name":       "users.full_name",
		"created_at": "users.created_at",
	}

	tests := []struct {
		name      string
		orders    []Order
		columnMap map[string]string
		wantQuery string
		wantErr   bool
	}{
		{
			name: "valid orders mapping",
			orders: []Order{
				{Column: "name", Method: OrderMethodASC},
				{Column: "created_at", Method: OrderMethodDESC},
			},
			columnMap: columnMap,
			wantQuery: "users.full_name ASC,users.created_at DESC",
			wantErr:   false,
		},
		{
			name:      "empty orders input",
			orders:    []Order{},
			columnMap: columnMap,
			wantQuery: "",
			wantErr:   false,
		},
		{
			name: "column not found in map",
			orders: []Order{
				{Column: "unknown_col", Method: OrderMethodASC},
			},
			columnMap: columnMap,
			wantQuery: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, err := OrdersQueryBuilder(tt.orders, tt.columnMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrdersQueryBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("OrdersQueryBuilder() = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}

func TestOrderMethod_String(t *testing.T) {
	if got := OrderMethodASC.String(); got != "ASC" {
		t.Errorf("OrderMethodASC.String() = %v, want %v", got, "ASC")
	}
	if got := OrderMethodDESC.String(); got != "DESC" {
		t.Errorf("OrderMethodDESC.String() = %v, want %v", got, "DESC")
	}
}
