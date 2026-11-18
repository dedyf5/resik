// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package status

import (
	"errors"
	"net/http"
	"testing"

	logCtx "github.com/dedyf5/resik/ctx/log"
	configEntity "github.com/dedyf5/resik/entities/config"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func TestGRPCErrorWithDetails(t *testing.T) {
	moduleConf := configEntity.Module{
		Name:    "RESTApp",
		NameKey: "rest-app-test",
	}
	dummyLog := logCtx.Get(configEntity.Log{}, moduleConf)

	t.Run("Nil error returns nil", func(t *testing.T) {
		res := GRPCErrorWithDetails(nil, dummyLog)
		if res != nil {
			t.Errorf("expecated nil, got %v", res)
		}
	})

	t.Run("Non-status error returns orignal error", func(t *testing.T) {
		oriErr := errors.New("standard error")
		res := GRPCErrorWithDetails(oriErr, dummyLog)
		if !errors.Is(res, oriErr) {
			t.Errorf("expected %v, got %v", oriErr, res)
		}
	})

	t.Run("Status without error condition or details returns original status error", func(t *testing.T) {
		statusErr := resPkg.NewStatusMessage(http.StatusOK, "Success", nil)
		res := GRPCErrorWithDetails(statusErr, dummyLog)
		if res == nil {
			t.Error("expected non-nil error, got nil")
		}
	})

	t.Run("Status error with details converts successfuly to gRPC status error", func(t *testing.T) {
		msg := "language is required"
		statusErr := resPkg.NewStatusBadRequest("en", "lang", msg, "REQUIRED", errors.New(msg))
		res := GRPCErrorWithDetails(statusErr, dummyLog)
		if res == nil {
			t.Error("expected non-nil gRPC error status, got nil")
		}
	})
}
