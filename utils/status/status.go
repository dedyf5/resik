// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package status

import (
	"errors"

	logCtx "github.com/dedyf5/resik/ctx/log"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func GRPCErrorWithDetails(err error, logCtx *logCtx.Log) error {
	if err == nil {
		return nil
	}

	statusPkg, ok := errors.AsType[*resPkg.Status](err)
	if !ok {
		return err
	}

	if !statusPkg.IsError() || len(statusPkg.Details) == 0 {
		return err
	}

	statusErr, unexpectedErr := statusPkg.GRPCStatus().WithDetails(statusPkg.Details...)
	if unexpectedErr != nil {
		logCtx.Error("failed to add details to status: " + unexpectedErr.Error())
		return err
	}

	return statusErr.Err()
}
