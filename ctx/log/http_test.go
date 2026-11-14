// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	configEntity "github.com/dedyf5/resik/entities/config"
	resPkg "github.com/dedyf5/resik/pkg/response"
	httpUtil "github.com/dedyf5/resik/utils/http"
	"go.uber.org/zap/zapcore"
)

func initLog() (loggerInst *Log, start time.Time, path string) {
	loggerInst = Get(configEntity.Log{File: ""}, configEntity.Module{NameKey: "http_module"})
	start = time.Now()
	path = "/user/profile?mode=1"
	return
}

func initHTTP(c context.Context, loggerInst *Log, start time.Time, path, contentType string, reqBody []byte) *HTTP {
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(c, http.MethodGet, path, nil)
	return NewHTTP(rec, c, loggerInst, start, req.Method, req.URL, contentType, req.UserAgent(), reqBody)
}

func TestNewHTTP(t *testing.T) {
	resetLog()
	loggerInst, start, path := initLog()
	c := context.Background()

	t.Run("Success HTTP request initialization", func(t *testing.T) {
		h := initHTTP(c, loggerInst, start, path, "application/json", nil)

		if h == nil {
			t.Fatalf("expected HTTP log instance, got nil")
		}
		if h.statusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, h.statusCode)
		}
		if expected := "mode=1"; h.url.RawQuery == "" || h.url.RawQuery != expected {
			t.Errorf("expected query %q, got %q", expected, h.url.RawQuery)
		}
	})
}

func TestHTTP_WriteHeader(t *testing.T) {
	resetLog()
	loggerInst, start, path := initLog()
	c := context.Background()

	h := initHTTP(c, loggerInst, start, path, "application/json", nil)

	if h == nil {
		t.Fatalf("expected HTTP log instance, got nil")
	}

	h.WriteHeader(http.StatusBadRequest)
	if h.statusCode != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, h.statusCode)
	}
}

func TestHTTP_Write(t *testing.T) {
	resetLog()
	loggerInst, start, path := initLog()

	t.Run("Info level log (2xx status) with caller from context", func(t *testing.T) {
		c := context.WithValue(context.Background(), KeyCallerHolderContext, &CallerHolder{
			Caller: "http_handler.go:30",
		})

		h := initHTTP(c, loggerInst, start, path, "application/json", nil)

		if h == nil {
			t.Fatalf("HTTP instance is nil")
		}

		_, err := h.Write([]byte(`{"status":"OK"}`))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("Warn level log (4xx status)", func(t *testing.T) {
		c := context.Background()
		h := initHTTP(c, loggerInst, start, path, "application/json", nil)
		h.statusCode = http.StatusBadRequest

		_, err := h.Write([]byte(`{"status":"Bad Request"}`))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("Error level log (5xx status) with stack trace", func(t *testing.T) {
		c := context.Background()

		msg := "Internal Server Error"
		status := resPkg.NewStatusError(http.StatusInternalServerError, errors.New(msg))

		h := initHTTP(c, loggerInst, start, path, "application/json", nil)
		h.statusCode = status.Code

		logResp := httpUtil.LoggerFromStatus(status)
		b, err := json.Marshal(logResp)
		if err != nil {
			t.Fatalf("failed json.Marshal: %v", err)
		}

		_, err = h.Write(b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestHTTP_MarshalLogObject(t *testing.T) {
	resetLog()
	loggerInst, start, path := initLog()

	t.Run("Marshal with multipart content type and response body", func(t *testing.T) {
		c := context.Background()
		h := initHTTP(c, loggerInst, start, path, "multipart/form-data; boundary=xyz", []byte("binary-data"))
		h.responseBody.WriteString("response output")

		enc := zapcore.NewMapObjectEncoder()
		err := h.MarshalLogObject(enc)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if expected := "[MULTIPART: binary data omitted]"; enc.Fields["request_body"] != expected {
			t.Errorf("expected request body %q, got %q", expected, enc.Fields["request_body"])
		}
	})

	t.Run("Marshal valid JSON request body with masking", func(t *testing.T) {
		c := context.Background()
		validJSONBody := []byte(`{"username": "dedy", "password": "secretpassword"}`)
		h := initHTTP(c, loggerInst, start, path, "application/json", validJSONBody)

		enc := zapcore.NewMapObjectEncoder()
		err := h.MarshalLogObject(enc)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		reqBodyStr, ok := enc.Fields["request_body"].(string)
		if !ok || reqBodyStr == "" {
			t.Fatal("expected request_body to be marshaled as string")
		}
	})

	t.Run("Marshal non-JSON request body exceeding 1000 characters (truncation)", func(t *testing.T) {
		c := context.Background()

		longPlainBody := make([]byte, 1050)
		for i := range longPlainBody {
			longPlainBody[i] = 'a'
		}

		h := initHTTP(c, loggerInst, start, path, "text/plain", longPlainBody)

		enc := zapcore.NewMapObjectEncoder()
		err := h.MarshalLogObject(enc)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		reqBodyStr, ok := enc.Fields["request_body"].(string)
		if !ok || len(reqBodyStr) <= 1000 {
			t.Fatal("expected truncated request body field")
		}
	})
}
