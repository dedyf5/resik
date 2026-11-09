// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import "testing"

func TestResponsePageSetup(t *testing.T) {
	int32Ptr := func(i int32) *int32 {
		return &i
	}

	tests := []struct {
		name        string
		total       int64
		limit       int32
		pageCurrent int32
		want        *ResponsePage
	}{
		{
			name:        "zero total items",
			total:       0,
			limit:       10,
			pageCurrent: 1,
			want: &ResponsePage{
				First:    1,
				Previous: nil,
				Current:  1,
				Next:     nil,
				Last:     1,
			},
		},
		{
			name:        "negative total items triggers last equals zero fallback",
			total:       -5,
			limit:       10,
			pageCurrent: 1,
			want: &ResponsePage{
				First:    1,
				Previous: nil,
				Current:  1,
				Next:     nil,
				Last:     1,
			},
		},
		{
			name:        "zero current page defaults to 1 (first page of multiple)",
			total:       25,
			limit:       10,
			pageCurrent: 0,
			want: &ResponsePage{
				First:    1,
				Previous: nil,
				Current:  1,
				Next:     int32Ptr(2),
				Last:     3,
			},
		},
		{
			name:        "middle page navigation (page 2 of 3)",
			total:       25,
			limit:       10,
			pageCurrent: 2,
			want: &ResponsePage{
				First:    1,
				Previous: int32Ptr(1),
				Current:  2,
				Next:     int32Ptr(3),
				Last:     3,
			},
		},
		{
			name:        "last page navigation (page 3 of 3)",
			total:       25,
			limit:       10,
			pageCurrent: 3,
			want: &ResponsePage{
				First:    1,
				Previous: int32Ptr(2),
				Current:  3,
				Next:     nil,
				Last:     3,
			},
		},
		{
			name:        "current page exceeds last page",
			total:       20,
			limit:       10,
			pageCurrent: 5,
			want: &ResponsePage{
				First:    1,
				Previous: int32Ptr(2),
				Current:  5,
				Next:     nil,
				Last:     2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResponsePageSetup(tt.total, tt.limit, tt.pageCurrent)
			if !compareResponsePages(got, tt.want) {
				t.Errorf("ResponsePageSetup() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestResponseMetaSetup(t *testing.T) {
	total := int64(50)
	limit := int32(10)
	pageCurrent := int32(2)

	meta := ResponseMetaSetup(total, limit, pageCurrent)

	if meta.GetTotal() != total {
		t.Errorf("ResponseMetaSetup().Total = %v, want %v", meta.GetTotal(), total)
	}
	if meta.GetLimit() != limit {
		t.Errorf("ResponseMetaSetup().Limit = %v, want %v", meta.GetLimit(), limit)
	}
	if meta.GetPage() == nil || meta.GetPage().GetCurrent() != pageCurrent {
		t.Errorf("ResponseMetaSetup().Page.Current = %v, want %v", meta.GetPage().GetCurrent(), pageCurrent)
	}
}

func compareResponsePages(a, b *ResponsePage) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.GetFirst() != b.GetFirst() || a.GetCurrent() != b.GetCurrent() || a.GetLast() != b.GetLast() {
		return false
	}

	if !compareInt32Ptr(a.Previous, b.Previous) || !compareInt32Ptr(a.Next, b.Next) {
		return false
	}

	return true
}

func compareInt32Ptr(a, b *int32) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
