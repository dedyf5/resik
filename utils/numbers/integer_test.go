// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package numbers

import "testing"

func TestSafeConvert(t *testing.T) {
	t.Run("Successful conversion from int to int64", func(t *testing.T) {
		val := 42
		res, err := SafeConvert[int64](val)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if res != 42 {
			t.Errorf("expected %d, got %d", 42, res)
		}
	})

	t.Run("Successful conversion from int to int32", func(t *testing.T) {
		val := 1000
		res, err := SafeConvert[int32](val)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if res != 1000 {
			t.Errorf("expected %d, got %d", 1000, res)
		}
	})

	t.Run("Conversion failure due to integer overflow", func(t *testing.T) {
		var val int64 = 3000000000
		res, err := SafeConvert[int32](val)
		if err == nil {
			t.Error("expected overflow error status, got nil")
		}
		if res != 0 {
			t.Errorf("expected zero result on error, got %v", res)
		}
	})
}
