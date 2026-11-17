// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package numbers

import (
	"net/http"

	"github.com/dedyf5/resik/pkg/numbers"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

// SafeConvert safely converts a numerical value from source type S to target type T.
//
// It the conversion fails due to overflow or incompatible values, it returns a zero value
// along with a standardized internal server error status wrapped in *resPkg.Status.
// Otherwise, it returns the converted value and a nil status error.
func SafeConvert[T numbers.Target, S numbers.Source](val S) (result T, err *resPkg.Status) {
	res, tmpErr := numbers.SafeConvert[T](val)
	if tmpErr != nil {
		return 0, resPkg.NewStatusError(http.StatusInternalServerError, tmpErr)
	}
	return res, nil
}
