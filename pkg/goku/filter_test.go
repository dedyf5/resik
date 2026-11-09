// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package goku

import "testing"

func TestFilter_FulltextSearch(t *testing.T) {
	tests := []struct {
		name   string
		search string
		want   string
	}{
		{
			name:   "standard multiple words",
			search: "apple banana",
			want:   "+apple* +banana*",
		},
		{
			name:   "search with special characters cleaned",
			search: "go+lang (backend)",
			want:   "+golang* +backend*",
		},
		{
			name:   "empty search string",
			search: "",
			want:   "",
		},
		{
			name:   "whitespace only search",
			search: "   ",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFilter(tt.search, 1, 10)
			if got := f.FulltextSeach(); got != tt.want {
				t.Errorf("Filter.FulltextSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Pagination(t *testing.T) {
	// Test default fallbacks
	fDefault := NewFilter("test", 0, 0)
	if got := fDefault.PageOrDefault(); got != int(PageDefault) {
		t.Errorf("PageOrDefault() = %v, want %v", got, PageDefault)
	}
	if got := fDefault.LimitOrDefault(); got != int(LimitDefault) {
		t.Errorf("LimitOrDefault() = %v, want %v", got, LimitDefault)
	}
	if got := fDefault.Offset(); got != 0 {
		t.Errorf("Offset() = %v, want %v", got, 0)
	}

	// Test custom values & offset calculation for page 2
	fCustom := NewFilter("test", 2, 15)
	if got := fCustom.PageOrDefault(); got != 2 {
		t.Errorf("PageOrDefault() = %v, want %v", got, 2)
	}
	if got := fCustom.LimitOrDefault(); got != 15 {
		t.Errorf("LimitOrDefault() = %v, want %v", got, 15)
	}
	if got := fCustom.Offset(); got != 15 {
		t.Errorf("Offset() = %v, want %v", got, 15)
	}

	// Test Raw pagination values
	rawDefault := fDefault.Raw()
	if rawDefault.PageOrDefault() != PageDefault || rawDefault.LimitOrDefault() != LimitDefault || rawDefault.Offset() != 0 {
		t.Errorf(
			"Raw pagination values mismatch: page=%d, limit=%d, offset=%d",
			rawDefault.PageOrDefault(), rawDefault.LimitOrDefault(), rawDefault.Offset(),
		)
	}
	rawCustom := fCustom.Raw()
	if rawCustom.PageOrDefault() != 2 || rawCustom.LimitOrDefault() != 15 || rawCustom.Offset() != 15 {
		t.Errorf(
			"Raw pagination values mismatch: page=%d, limit=%d, offset=%d",
			rawCustom.PageOrDefault(), rawCustom.LimitOrDefault(), rawCustom.Offset(),
		)
	}
}
