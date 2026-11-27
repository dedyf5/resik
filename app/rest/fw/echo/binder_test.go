// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/labstack/echo/v5"
	"golang.org/x/text/language"
)

type bindPayload struct {
	ID       int      `query:"id" param:"id" form:"id" json:"id"`
	Name     string   `query:"name" form:"name" json:"name"`
	Price    float64  `query:"price" form:"price" json:"price"`
	Active   bool     `query:"active" form:"active" json:"active"`
	Tags     []string `query:"tags" form:"tags" json:"tags"`
	PtrCount *int     `query:"ptr_count" form:"ptr_count" json:"ptr_count"`
}

func newBinderTestContext(t *testing.T, method, target, body string, withLang bool) (*echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	e := echo.New()
	c := context.Background()
	var req *http.Request
	if body != "" {
		req = httptest.NewRequestWithContext(c, method, target, bytes.NewBufferString(body))
	} else {
		req = httptest.NewRequestWithContext(c, method, target, nil)
	}

	if withLang {
		lang := langCtx.NewLang(language.English, &language.Indonesian, "id")
		c = context.WithValue(req.Context(), langCtx.ContextKey, lang)
		req = req.WithContext(c)
	}

	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestNewBinder(t *testing.T) {
	b1 := NewBinder()
	b2 := NewBinder()

	if b1 == nil {
		t.Fatal("expected non-nil Binder")
	}
	if b1 != b2 {
		t.Error("expected NewBinder to return singleton instance")
	}
}

func TestBinder_BindHelpers(t *testing.T) {
	b := NewBinder()

	t.Run("BindHeaders", func(t *testing.T) {
		headerVal := "header-val"
		c, _ := newBinderTestContext(t, http.MethodGet, "/items", "", true)
		c.Request().Header.Set("X-Custom-Name", headerVal)

		var dst struct {
			CustomName string `header:"X-Custom-Name"`
		}
		if err := b.BindHeaders(c, &dst); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dst.CustomName != headerVal {
			t.Errorf("CustomName = %q, want %q", dst.CustomName, headerVal)
		}
	})

	t.Run("BindBody", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodPost, "/items", `{"name":"item-a","id":1}`, true)
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		var dst bindPayload
		if err := b.BindBody(c, &dst); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dst.Name != "item-a" {
			t.Errorf("Name = %q, want %q", dst.Name, "item-a")
		}
		if dst.ID != 1 {
			t.Errorf("ID = %d, want %d", dst.ID, 1)
		}
	})

	t.Run("BindQueryParams", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodGet, "/items?name=query-val&id=7", "", true)

		var dst bindPayload
		if err := b.BindQueryParams(c, &dst); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dst.Name != "query-val" {
			t.Errorf("Name = %q, want %q", dst.Name, "query-val")
		}
		if dst.ID != 7 {
			t.Errorf("ID = %d, want %d", dst.ID, 7)
		}
	})

	t.Run("BindPathParams", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodGet, "/items/42", "", true)
		c.SetPath("/items/:id")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: "42"}})

		var dst bindPayload
		if err := b.BindPathParams(c, &dst); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dst.ID != 42 {
			t.Errorf("ID = %d, want %d", dst.ID, 42)
		}
	})
}

func TestBinder_ParamValidator(t *testing.T) {
	b := NewBinder()

	t.Run("Return nil when payload is nil or non-struct", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodGet, "/items", "", true)

		if err := b.ParamValidator(c, nil); err != nil {
			t.Fatalf("expected nil for nil target, got %v", err)
		}

		strVal := "not-a-struct"
		if err := b.ParamValidator(c, &strVal); err != nil {
			t.Fatalf("expected nil for pointer to non-struct, got %v", err)
		}
	})

	t.Run("Return error when query lang is unavailable", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodGet, "/items?lang=xx-INVALID", "", true)

		var dst bindPayload
		if err := b.ParamValidator(c, &dst); err == nil {
			t.Fatal("expected error for unavailable query language, got nil")
		}
	})

	t.Run("Return error when language context is missing", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodGet, "/items?id=1", "", false)

		var dst bindPayload
		if err := b.ParamValidator(c, &dst); err == nil {
			t.Fatal("expected error when language context is missing, got nil")
		}
	})

	t.Run("Pass with valid query, param, and form field types", func(t *testing.T) {
		form := url.Values{}
		form.Set("active", "true")
		form.Set("ptr_count", "-5")

		c, _ := newBinderTestContext(
			t,
			http.MethodPost,
			"/items/99?name=ok&price=12.5&tags=a&tags=b",
			form.Encode(),
			true,
		)
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		c.SetPath("/items/:id")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: "99"}})

		var dst bindPayload
		if err := b.ParamValidator(c, &dst); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Fail when query int has invalid type", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodGet, "/items?id=abc", "", true)

		var dst bindPayload
		if err := b.ParamValidator(c, &dst); err == nil {
			t.Fatal("expected validation error for invalid int query param, got nil")
		}
	})

	t.Run("Fail when path param has invalid type", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodGet, "/items/not-int", "", true)
		c.SetPath("/items/:id")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: "not-int"}})

		var dst bindPayload
		if err := b.ParamValidator(c, &dst); err == nil {
			t.Fatal("expected validation error for invalid path param, got nil")
		}
	})

	t.Run("Fail when form bool has invalid type", func(t *testing.T) {
		form := url.Values{}
		form.Set("active", "not-bool")

		c, _ := newBinderTestContext(t, http.MethodPost, "/items", form.Encode(), true)
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

		var dst bindPayload
		if err := b.ParamValidator(c, &dst); err == nil {
			t.Fatal("expected validation error for invalid form bool, got nil")
		}
	})
}

func TestBinder_JSONErrorFormatter(t *testing.T) {
	b := NewBinder()

	t.Run("Return error when language context is missing", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodPost, "/items", "", false)

		err := b.JSONErrorFormatter(c, errors.New("some error"))
		if err == nil {
			t.Fatal("expected error when language context is missing, got nil")
		}
	})

	t.Run("Format *json.UnmarshalTypeError", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodPost, "/items", "", true)

		unmarshalErr := &json.UnmarshalTypeError{
			Value: "string",
			Type:  reflect.TypeOf(int(0)),
			Field: "id",
		}
		err := b.JSONErrorFormatter(c, unmarshalErr)
		if err == nil {
			t.Fatal("expected formatted error for UnmarshalTypeError, got nil")
		}
	})

	t.Run("Format error matching field, expecated, got pattern", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodPost, "/items", "", true)

		formattedErr := errors.New("failed to bind: field=age,expected=int,got=abc,")
		err := b.JSONErrorFormatter(c, formattedErr)
		if err == nil {
			t.Fatalf("expected formatted error for field/expected/got string, got nil")
		}
	})

	t.Run("Return original error when pattern does not match", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodPost, "/items", "", true)

		rawErr := errors.New("unexpected binding failure")
		err := b.JSONErrorFormatter(c, rawErr)
		if !errors.Is(err, rawErr) {
			t.Errorf("error = %v, want %v", err, rawErr)
		}
	})
}

func TestBinder_Bind(t *testing.T) {
	b := NewBinder()

	t.Run("Fail when ParamValidator rejects parameter", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodGet, "/items?id=invalid", "", true)

		var dst bindPayload
		if err := b.Bind(c, &dst); err == nil {
			t.Fatal("expected error from ParamValidator, got nil")
		}
	})

	t.Run("Fail when default binder fails and format JSON error", func(t *testing.T) {
		body := `{"name":"ok","id":"not-an-int"}`
		c, _ := newBinderTestContext(t, http.MethodPost, "/items", strings.TrimSpace(body), true)
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		var dst bindPayload
		if err := b.Bind(c, &dst); err == nil {
			t.Fatalf("expected error when binding invalid JSON field type, got nil")
		}
	})

	t.Run("Succeed when parameters and body are valid", func(t *testing.T) {
		c, _ := newBinderTestContext(t, http.MethodPost, "/items", `{"id":15}`, true)
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		var dst bindPayload
		if err := b.Bind(c, &dst); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dst.ID != 15 {
			t.Errorf("dst.ID = %d, want %d", dst.ID, 15)
		}
	})
}
