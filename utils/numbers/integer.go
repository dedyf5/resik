// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package numbers

import (
	"net/http"

	"github.com/dedyf5/resik/pkg/numbers"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func SafeConvert[T numbers.Target, S numbers.Source](val S) (result T, err *resPkg.Status) {
	res, tmpErr := numbers.SafeConvert[T](val)
	if tmpErr != nil {
		return 0, resPkg.NewStatusError(http.StatusInternalServerError, tmpErr)
	}
	return res, nil
}
