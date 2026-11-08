// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package collection

import (
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		mapper func(int) string
		want   []string
	}{
		{
			name:  "transform int slice to string slice",
			input: []int{1, 2, 3},
			mapper: func(i int) string {
				return strconv.Itoa(i * 10)
			},
			want: []string{"10", "20", "30"},
		},
		{
			name:  "empty slice input",
			input: []int{},
			mapper: func(i int) string {
				return strconv.Itoa(i * 10)
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.input, tt.mapper)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	const (
		tech    = "Tech"
		finance = "Finance"
		health  = "Health"
	)

	type Sample struct {
		ID       int
		Category string
	}

	tests := []struct {
		name      string
		items     []Sample
		extractor func(Sample) string
		want      []string
	}{
		{
			name: "extract unique string keys preserving order",
			items: []Sample{
				{ID: 1, Category: tech},
				{ID: 2, Category: finance},
				{ID: 3, Category: tech},
				{ID: 4, Category: health},
				{ID: 5, Category: finance},
			},
			extractor: func(s Sample) string {
				return s.Category
			},
			want: []string{tech, finance, health},
		},
		{
			name:      "empty items slice",
			items:     []Sample{},
			extractor: func(s Sample) string { return s.Category },
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.items, tt.extractor)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}
