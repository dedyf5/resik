// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"errors"
	"net/http"
	"testing"

	resPkg "github.com/dedyf5/resik/pkg/response"
)

const (
	success = "Success"
	okCode  = "200.1"
)

func TestResponseSuccessAuto(t *testing.T) {
	t.Run("Success response without meta", func(t *testing.T) {
		status := &resPkg.Status{
			Code:    http.StatusOK,
			Message: success,
			Data:    "test-data",
		}

		res := ResponseSuccessAuto(status)
		if res.Status.Code != okCode {
			t.Errorf("expected code %q, got %q", okCode, res.Status.Code)
		}
		if res.Status.Message != success {
			t.Errorf("expected message %q, got %q", success, res.Status.Message)
		}
		if res.Data != "test-data" {
			t.Errorf("expected data %q, got %q", "test-data", res.Data)
		}
		if res.Meta != nil {
			t.Errorf("expected nil meta, got %v", res.Meta)
		}
	})

	t.Run("Success response with meta", func(t *testing.T) {
		status := &resPkg.Status{
			Code:    http.StatusOK,
			Message: success,
			Meta: &resPkg.Meta{
				Total:       100,
				Limit:       10,
				PageCurrent: 1,
			},
		}

		res := ResponseSuccessAuto(status)
		if res.Meta == nil {
			t.Fatalf("expected non-nil meta, got nil")
		}
	})
}

func TestResponseErrorAuto(t *testing.T) {
	t.Run("Error response without details", func(t *testing.T) {
		status := &resPkg.Status{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}

		res := ResponseErrorAuto(status)
		if expected := "500.1"; res.Status.Code != expected {
			t.Errorf("expected code %q, got %q", expected, res.Status.Code)
		}
		if expected := "Internal Server Error"; res.Status.Message != expected {
			t.Errorf("expected message %q, got %q", expected, res.Status.Message)
		}
		if res.Status.Details != nil {
			t.Errorf("expected nil details, got %v", res.Status.Details)
		}
	})

	t.Run("Error response with details", func(t *testing.T) {
		msg := "language is required"
		status := resPkg.NewStatusBadRequest("en", "lang", msg, "REQUIRED", errors.New(msg))

		res := ResponseErrorAuto(status)
		if expected := "400.1"; res.Status.Code != expected {
			t.Errorf("expected code %q, got %q", expected, res.Status.Code)
		}
		if res.Status.Details == nil {
			t.Error("expoected non-nil details, got nil")
		}
	})

	t.Run("Error response with data", func(t *testing.T) {
		status := &resPkg.Status{
			Code: http.StatusUnauthorized,
			Data: "error-detail-data",
		}

		res := ResponseErrorAuto(status)
		if expected := "error-detail-data"; res.Data != expected {
			t.Errorf("expected error data %q, got %q", expected, res.Data)
		}
	})
}

func TestResponseFromStatus(t *testing.T) {
	t.Run("Route to success response", func(t *testing.T) {
		status := resPkg.NewStatusCode(http.StatusOK)
		res := ResponseFromStatus(status)
		if res.Status.Code != okCode {
			t.Errorf("expected success code %q, got %q", okCode, res.Status.Code)
		}
	})

	t.Run("Route to error response", func(t *testing.T) {
		status := resPkg.NewStatusError(http.StatusInternalServerError, errors.New("internal error"))
		res := ResponseFromStatus(status)
		if expected := "500.1"; res.Status.Code != expected {
			t.Errorf("expected error code %q, got %q", expected, res.Status.Code)
		}
	})
}

func TestLoggerFromStatus(t *testing.T) {
	status := &resPkg.Status{
		Code:   http.StatusOK,
		Caller: "test-caller",
	}

	log := LoggerFromStatus(status)
	if expected := "test-caller"; log.Caller != expected {
		t.Errorf("expected caller %q, got %q", expected, log.Caller)
	}
	if log.Response.Status.Code != okCode {
		t.Errorf("expected response code %q, got %q", okCode, log.Response.Status.Code)
	}
}

func TestLoggerErrorAuto(t *testing.T) {
	status := &resPkg.Status{
		Code:   http.StatusBadRequest,
		Caller: "error-caller",
	}

	log := LoggerErrorAuto(status)
	if expected := "error-caller"; log.Caller != expected {
		t.Errorf("expected caller %q, got %q", expected, log.Caller)
	}
	if expected := "400.1"; log.Response.Status.Code != expected {
		t.Errorf("expected response code %q, got %q", expected, log.Response.Status.Code)
	}
}
