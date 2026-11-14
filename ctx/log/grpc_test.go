// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	commonCore "github.com/dedyf5/resik/core/common"
	commonEntity "github.com/dedyf5/resik/entities/common"
	configEntity "github.com/dedyf5/resik/entities/config"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"go.uber.org/zap/zapcore"
)

type mockGRPCStatusResponse struct {
	status *commonCore.Status
}

func (m *mockGRPCStatusResponse) GetStatus() *commonCore.Status {
	return m.status
}

func initGRPC() (loggerInst *Log, start time.Time, path string) {
	loggerInst = Get(configEntity.Log{File: ""}, configEntity.Module{NameKey: "grpc_module"})
	start = time.Now()
	path = "/user.Service/Profile"
	return
}

func TestNewGRPC(t *testing.T) {
	resetLog()

	loggerInst, start, path := initGRPC()
	c := context.Background()

	t.Run("Success response without error", func(t *testing.T) {
		customStatus := commonCore.Status{Code: "200.1", Message: "custom success message"}

		resp := &mockGRPCStatusResponse{status: &customStatus}
		g := NewGRPC(c, loggerInst, start, path, commonEntity.Map{"id": 1}, resp, nil)
		if g == nil {
			t.Fatal("expecated GRPC instance, got nil")
		}
		if g.status.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, g.status.Code)
		}
		if g.status.Message != customStatus.GetMessage() {
			t.Errorf("expected status message %q, got %q", customStatus.GetMessage(), g.status.Message)
		}
	})

	t.Run("Error response with resPkg.Status", func(t *testing.T) {
		statusErr := resPkg.NewStatusError(http.StatusBadRequest, errors.New("invalid input"))
		g := NewGRPC(c, loggerInst, start, path, nil, nil, statusErr)
		if g.status.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, g.status.Code)
		}
	})

	t.Run("Error response with generic error", func(t *testing.T) {
		genericErr := errors.New("generic internal error")
		g := NewGRPC(c, loggerInst, start, path, nil, nil, genericErr)
		if g.status.Code != http.StatusOK {
			t.Errorf("expected default status code %d, got %d", http.StatusOK, g.status.Code)
		}
	})
}

func TestGRPC_Write(t *testing.T) {
	resetLog()

	loggerInst, start, path := initGRPC()
	loggerInst.CorrelationID = "grpc-corr-123"

	t.Run("Info level log with caller from context", func(t *testing.T) {
		c := context.WithValue(context.Background(), KeyCallerHolderContext, &CallerHolder{
			Caller: "user_handler.go:45",
		})
		g := NewGRPC(c, loggerInst, start, path,
			commonEntity.Map{"password": "secret"},
			commonEntity.Map{"ok": true},
			nil,
		)
		g.Write()
	})

	t.Run("Warn level log (4xx status) with caller override on status", func(t *testing.T) {
		statusErr := resPkg.NewStatusError(http.StatusNotFound, errors.New("user not found"))
		statusErr.Caller = "user_repo.go:88"

		g := NewGRPC(context.Background(), loggerInst, start, path, nil, nil, statusErr)
		g.Write()
	})

	t.Run("Error level log (5xx status) with stack trace", func(t *testing.T) {
		statusErr := resPkg.NewStatusError(http.StatusInternalServerError, errors.New("database timeout"))
		statusErr.StackTrace = []string{"goroutine 1 [running]:", "main.main()"}

		g := NewGRPC(context.Background(), loggerInst, start, path, nil, nil, statusErr)
		g.Write()
	})
}

func TestGRPC_MarshalLogObject(t *testing.T) {
	resetLog()

	loggerInst, start, path := initGRPC()
	enc := zapcore.NewMapObjectEncoder()

	t.Run("Success payload with string response", func(t *testing.T) {
		resBody := "plain string response"
		g := NewGRPC(context.Background(), loggerInst, start, path,
			commonEntity.Map{"token": "abc"},
			resBody,
			nil,
		)

		if err := g.MarshalLogObject(enc); err != nil {
			t.Fatalf("unexpected error in MarshalLogObject: %v", err)
		}

		if enc.Fields["module"] != "grpc_module" {
			t.Errorf("expected module %q, got %q", "grpc_module", enc.Fields["module"])
		}
		if enc.Fields["response_body"] != resBody {
			t.Errorf("expected response body %q, got %q", resBody, enc.Fields["response_body"])
		}
	})

	t.Run("Unmarshalable requestBody triggers request_body_error", func(t *testing.T) {
		encErr := zapcore.NewMapObjectEncoder()
		invalidReq := make(chan int)

		g := NewGRPC(context.Background(), loggerInst, start, path,
			invalidReq,
			commonEntity.Map{"status": "ok"},
			nil,
		)

		if err := g.MarshalLogObject(encErr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, exists := encErr.Fields["request_body_error"]; !exists {
			t.Errorf("expected 'request_body_error' field to be populated")
		}
	})

	t.Run("Error status logs message as response_body", func(t *testing.T) {
		encStatusErr := zapcore.NewMapObjectEncoder()
		msg := "Language is required"
		statusErr := resPkg.NewStatusBadRequest("en", "lang", msg, "REQUIRED", errors.New(msg))
		g := NewGRPC(context.Background(), loggerInst, start, path, nil, nil, statusErr)

		if err := g.MarshalLogObject(encStatusErr); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := encStatusErr.Fields["response_body"]; got != statusErr.MessageOrDefault() {
			t.Errorf("expected response body %q, got %q", statusErr.MessageOrDefault(), got)
		}
	})
}
