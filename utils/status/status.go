// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package status

import (
	"errors"

	logCtx "github.com/dedyf5/resik/ctx/log"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func GRPCErrorWithDetails(err error, logCtx *logCtx.Log) error {
	if err == nil {
		return nil
	}

	statusPkg, ok := errors.AsType[*resPkg.Status](err)
	if !ok {
		return err
	}

	if !statusPkg.IsError() || len(statusPkg.Detail) == 0 {
		return err
	}

	fieldViolations := make([]*errdetails.BadRequest_FieldViolation, 0, len(statusPkg.Detail))
	for field, description := range statusPkg.Detail {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: description,
		})
	}

	badRequest := errdetails.BadRequest{
		FieldViolations: fieldViolations,
	}

	statusErr, unexpectedErr := statusPkg.GRPCStatus().WithDetails(&badRequest)
	if unexpectedErr != nil {
		logCtx.Error("failed to add details to status: " + unexpectedErr.Error())
		return err
	}

	return statusErr.Err()
}
