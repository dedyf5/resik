// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package status

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

func CodePlus(code codes.Code) string {
	return fmt.Sprintf("%d.1", code)
}
