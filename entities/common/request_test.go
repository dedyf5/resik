// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package common

import (
	"reflect"
	"testing"
)

func TestRequest_LangAvailable(t *testing.T) {
	req := &Request{}
	got := req.LangAvailable()
	want := []string{"en", "id", "ja"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("LangAvailable() = %v, want %v", got, want)
	}
}
