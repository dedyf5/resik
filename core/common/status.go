// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package common

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

type statusGetter interface {
	GetStatus() *Status
}

// Extract pulls the Status message from an arbitrary object, typically a gRPC response.
//
// Returns nil if the object is nil, or if it does not satisfy the status getter interface.
func StatusExtract(param any) *Status {
	if param == nil {
		return nil
	}

	if s, ok := param.(statusGetter); ok {
		return s.GetStatus()
	}

	return nil
}

func StatusCodePlus(code codes.Code) string {
	return fmt.Sprintf("%d.1", code)
}
