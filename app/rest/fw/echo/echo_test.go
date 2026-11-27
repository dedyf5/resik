// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	resPkg "github.com/dedyf5/resik/pkg/response"
	mockValidator "github.com/dedyf5/resik/utils/validator/mock"
	"github.com/labstack/echo/v5"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
)

type simplePayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type customStatusCoderErr struct {
	code int
	msg  string
}

func (e *customStatusCoderErr) Error() string {
	return e.msg
}

func (e *customStatusCoderErr) StatusCode() int {
	return e.code
}

func newEchoContext(t *testing.T, method, target, body string, withLang bool) (*echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	e := echo.New()
	c := context.Background()
	var req *http.Request
	if body != "" {
		req = httptest.NewRequestWithContext(c, method, target, bytes.NewBufferString(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequestWithContext(c, method, target, nil)
	}

	if withLang {
		lang := langCtx.NewLang(language.English, nil, "en")
		c = context.WithValue(req.Context(), langCtx.ContextKey, lang)
		req = req.WithContext(c)
	}

	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func getMockValidator(t *testing.T) (mockVal *mockValidator.MockIValidate) {
	ctrl := gomock.NewController(t)
	return mockValidator.NewMockIValidate(ctrl)
}

func TestEcho_StructValidator(t *testing.T) {
	t.Run("Fail when reqest binding fails", func(t *testing.T) {
		mockVal := getMockValidator(t)

		ic := New(mockVal)
		c, _ := newEchoContext(t, http.MethodPost, "/users", `{"age":"not-an-int"}`, true)

		var payload simplePayload
		err := ic.StructValidator(c, &payload)
		if err == nil {
			t.Fatal("expected bind error, got nil")
		}
	})

	t.Run("Fail when language context is missing", func(t *testing.T) {
		mockVal := getMockValidator(t)

		ic := New(mockVal)
		c, _ := newEchoContext(t, http.MethodPost, "/users", `{"name":"budi","age":25}`, false)

		var payload simplePayload
		err := ic.StructValidator(c, &payload)
		if err == nil {
			t.Fatal("expected error when language context is missing, got nil")
		}
	})

	t.Run("Fail when struct validation returns error", func(t *testing.T) {
		mockVal := getMockValidator(t)

		wantStatus := resPkg.NewStatusCode(http.StatusBadRequest)
		mockVal.EXPECT().
			Struct(gomock.Any(), gomock.Any()).
			Return(wantStatus)

		ic := New(mockVal)
		c, _ := newEchoContext(t, http.MethodPost, "/users", `{"name":"budi","age":25}`, true)

		var payload simplePayload
		err := ic.StructValidator(c, &payload)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("Succeed when bind and struct validation pass", func(t *testing.T) {
		mockVal := getMockValidator(t)

		mockVal.EXPECT().Struct(gomock.Any(), gomock.Any()).Return(nil)

		ic := New(mockVal)
		c, _ := newEchoContext(t, http.MethodPost, "/users", `{"name":"budi","age":25}`, true)

		var payload simplePayload
		err := ic.StructValidator(c, &payload)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if payload.Name != "budi" {
			t.Fatalf("expected payload name to be 'budi', got '%s'", payload.Name)
		}
		if payload.Age != 25 {
			t.Fatalf("expected payload age to be 25, got %d", payload.Age)
		}
	})
}

func TestHTTPErrorHandler(t *testing.T) {
	t.Run("Skip writing response when response is already commited", func(t *testing.T) {
		bodyString := "already-sent"
		c, rec := newEchoContext(t, http.MethodGet, "/test", "", false)
		rec.WriteHeader(http.StatusAccepted)
		_, _ = rec.WriteString(bodyString)

		resp, err := echo.UnwrapResponse(c.Response())
		if err == nil && resp != nil {
			resp.Committed = true
		}

		HTTPErrorHandler(c, errors.New("ignored error"))

		if rec.Code != http.StatusAccepted {
			t.Fatalf("status code = %d, want %d", rec.Code, http.StatusAccepted)
		}
		if rec.Body.String() != bodyString {
			t.Fatalf("body = %q, want %q", rec.Body.String(), bodyString)
		}
	})

	t.Run("Handle *resPkg.Status error", func(t *testing.T) {
		c, rec := newEchoContext(t, http.MethodGet, "/test", "", false)
		statusErr := resPkg.NewStatusCode(http.StatusBadRequest)

		HTTPErrorHandler(c, statusErr)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status code = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("Handle echo.HTTPStatusCoder error", func(t *testing.T) {
		c, rec := newEchoContext(t, http.MethodGet, "/test", "", false)
		coderErr := &customStatusCoderErr{code: http.StatusNotFound, msg: "not found"}

		HTTPErrorHandler(c, coderErr)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status code = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("Handle generic error with HTTP 500", func(t *testing.T) {
		c, rec := newEchoContext(t, http.MethodGet, "/test", "", false)

		HTTPErrorHandler(c, errors.New("unexpected failure"))

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status code = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})
}
