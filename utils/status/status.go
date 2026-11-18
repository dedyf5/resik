// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package status provides utility functions for handling and converting
// application status structures, gRPC error mappings, and error details enrichment.
package status

import (
	"errors"

	logCtx "github.com/dedyf5/resik/ctx/log"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

// GRPCErrorWithDetails converts a standard error into a gRPC status error enriched with detail payloads
// if the underlying error is a valid *resPkg.Status containing error details.
//
// if the error is nil, not *resPkg.Status, not an error status condition, or lacks details,
// it returns the original error unmodified. If appending details fails, it logs the unexpected error
// using the provided log context and falls back to returning the original error.
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
