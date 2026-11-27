// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtCtx "github.com/dedyf5/resik/ctx/jwt"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	mockIdentity "github.com/dedyf5/resik/internal/identity/mock"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/ratelimit"
	mockRatelimit "github.com/dedyf5/resik/utils/ratelimit/mock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
)

func newMiddlewareTestContext(t *testing.T, method, target, body string, withLang bool) (*echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	e := echo.New()
	c := context.Background()
	var req *http.Request
	if body != "" {
		req = httptest.NewRequestWithContext(c, method, target, bytes.NewBufferString(body))
	} else {
		req = httptest.NewRequestWithContext(c, method, target, nil)
	}
	req.RequestURI = target

	if withLang {
		lang := langCtx.NewLang(language.English, &language.Indonesian, "id")
		c = context.WithValue(req.Context(), langCtx.ContextKey, lang)
		req = req.WithContext(c)
	}

	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func signTestJWT(t *testing.T, secret string, exp time.Time) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user-123",
		"exp": exp.Unix(),
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign test JWT: %v", err)
	}

	return signed
}

func newLimiter(t *testing.T) *mockRatelimit.MockLimiter {
	ctrl := gomock.NewController(t)
	return mockRatelimit.NewMockLimiter(ctrl)
}

func TestLangMiddleware(t *testing.T) {
	mw := LangMiddleware(language.Indonesian)

	// t.Run("Set default language when query and header are empty", func(t *testing.T) {
	// 	c, rec := newMiddlewareTestContext(t, http.MethodGet, "/ping", "", false)

	// 	var resolved *langCtx.Lang
	// 	var fromCtxErr error
	// 	handler := mw(func(ctx *echo.Context) error {
	// 		resolved, fromCtxErr = langCtx.FromContext(ctx.Request().Context())
	// 		if fromCtxErr != nil {
	// 			return fromCtxErr
	// 		}
	// 		return ctx.String(http.StatusOK, "OK")
	// 	})

	// 	if err := handler(c); err != nil {
	// 		t.Fatalf("unexpected error: %v", err)
	// 	}
	// 	if rec.Code != http.StatusOK {
	// 		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	// 	}
	// 	if resolved == nil {
	// 		t.Fatal("expected resolved lang in context, got nil")
	// 	}
	// })

	t.Run("Resolve language from query param when available", func(t *testing.T) {
		target := "/ping?" + langCtx.ContextKey.String() + "=en"
		c, _ := newMiddlewareTestContext(t, http.MethodGet, target, "", false)

		handler := mw(func(ctx *echo.Context) error {
			resolved, err := langCtx.FromContext(ctx.Request().Context())
			if err != nil {
				return err
			}
			if resolved.LanguageReqOrDefault() != language.English {
				t.Errorf("got lang = %v, want %v", resolved.LanguageReqOrDefault(), language.English)
			}
			return ctx.NoContent(http.StatusOK)
		})

		if err := handler(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	// t.Run("Fallback when query language is unavailable", func(t *testing.T) {
	// 	target := "/ping?" + langCtx.ContextKey.String() + "=xx-unknown"
	// 	c, _ := newMiddlewareTestContext(t, http.MethodGet, target, "", false)
	// 	c.Request().Header.Set("Accept-Language", "en-US")

	// 	handler := mw(func(ctx *echo.Context) error {
	// 		if _, err := langCtx.FromContext(ctx.Request().Context()); err != nil {
	// 			return err
	// 		}
	// 		return ctx.NoContent(http.StatusOK)
	// 	})

	// 	if err := handler(c); err != nil {
	// 		t.Fatalf("unexpected error: %v", err)
	// 	}
	// })
}

func TestLoggerAndResponseFormatterMiddleware(t *testing.T) {
	baseLog := &logCtx.Log{}
	mw := LoggerAndResponseFormatterMiddleware(baseLog)

	t.Run("Bypass documentation prefix", func(t *testing.T) {
		c, rec := newMiddlewareTestContext(t, http.MethodGet, DocPrefix+"/index.html", "", false)

		called := false
		swagger := "swagger"
		handler := mw(func(ctx *echo.Context) error {
			called = true
			return ctx.String(http.StatusOK, swagger)
		})

		if err := handler(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Error("expected next handler to be called")
		}
		if rec.Body.String() != swagger {
			t.Errorf("body = %q, want %q", rec.Body.String(), swagger)
		}
	})

	t.Run("Read and restore request body and set context values", func(t *testing.T) {
		bodyStr := `{"item":"book"}`
		c, _ := newMiddlewareTestContext(t, http.MethodPost, "/orders", bodyStr, false)
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		var bodyAfter string
		var hasCallerHolder bool
		handler := mw(func(ctx *echo.Context) error {
			b, err := io.ReadAll(ctx.Request().Body)
			if err != nil {
				return err
			}
			bodyAfter = string(b)
			_, hasCallerHolder = ctx.Request().Context().Value(logCtx.KeyCallerHolderContext).(*logCtx.CallerHolder)
			return ctx.NoContent(http.StatusCreated)
		})

		if err := handler(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if bodyAfter != bodyStr {
			t.Errorf("bodyAfter = %q, want %q", bodyAfter, bodyStr)
		}
		if !hasCallerHolder {
			t.Error("expected CallerHolder in request context")
		}
	})

	t.Run("Skip body read for multipart/form-data", func(t *testing.T) {
		c, _ := newMiddlewareTestContext(t, http.MethodPost, "/upload", "raw-multipart", false)
		c.Request().Header.Set(echo.HeaderContentType, "multipart/form-data; boundary=---123")

		handler := mw(func(ctx *echo.Context) error {
			return ctx.NoContent(http.StatusOK)
		})

		if err := handler(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestValidateTokenMiddleware(t *testing.T) {
	const secret = "test-secret-key"
	mw := ValidateTokenMiddleware(secret)

	t.Run("Allow request when token is valid", func(t *testing.T) {
		c, _ := newMiddlewareTestContext(t, http.MethodGet, "/profile", "", true)
		token := signTestJWT(t, secret, time.Now().Add(time.Hour))
		c.Request().Header.Set("Authorization", "Bearer "+token)

		called := false
		handler := mw(func(ctx *echo.Context) error {
			called = true
			return ctx.NoContent(http.StatusOK)
		})

		if err := handler(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Error("expected next handler to be called for valid token")
		}
	})

	t.Run("Return formatted error when token is missing and lang context exists", func(t *testing.T) {
		c, _ := newMiddlewareTestContext(t, http.MethodGet, "/profile", "", true)

		handler := mw(func(ctx *echo.Context) error {
			t.Error("handler should not be called when token is missing")
			return nil
		})

		if err := handler(c); err == nil {
			t.Fatal("expected error for missing token, got nil")
		}
	})

	t.Run("Return error when token is missing and lang context is absent", func(t *testing.T) {
		c, _ := newMiddlewareTestContext(t, http.MethodGet, "/profile", "", false)

		handler := mw(func(ctx *echo.Context) error {
			t.Error("handler should not be called when token is missing")
			return nil
		})

		if err := handler(c); err == nil {
			t.Errorf("expected error when token and lang context are missing, got nil")
		}
	})
}

func TestJWTMiddleware(t *testing.T) {
	const secret = "test-jwt-secret"

	t.Run("Pass through when Authorization header is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockResolver := mockIdentity.NewMockIdentityResolver(ctrl)

		mw := JWTMiddleware(secret, mockResolver, language.English)
		c, _ := newMiddlewareTestContext(t, http.MethodGet, "/public", "", true)

		called := false
		handler := mw(func(ctx *echo.Context) error {
			called = true
			return ctx.NoContent(http.StatusOK)
		})

		if err := handler(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Error("expected next handler to be called without Authorization header")
		}
	})

	t.Run("Return error when token is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockResolver := mockIdentity.NewMockIdentityResolver(ctrl)

		mw := JWTMiddleware(secret, mockResolver, language.English)
		c, _ := newMiddlewareTestContext(t, http.MethodGet, "/private", "", false)
		c.Request().Header.Set("Authorization", "Bearer invalid-token")

		handler := mw(func(ctx *echo.Context) error {
			t.Error("handler should not be called when token parsing fails")
			return nil
		})

		if err := handler(c); err == nil {
			t.Fatal("expected error for invalid token, got nil")
		}
	})

	t.Run("Set auth claims in context when token is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockResolver := mockIdentity.NewMockIdentityResolver(ctrl)

		mockResolver.EXPECT().
			GetTenantOrganizationIDs(gomock.Any(), uint64(0)).
			Return([]uint64{}, nil)
		mockResolver.EXPECT().
			GetTenantBranchIDs(gomock.Any(), uint64(0)).
			Return([]uint64{}, nil)

		mw := JWTMiddleware(secret, mockResolver, language.English)
		c, _ := newMiddlewareTestContext(t, http.MethodGet, "/private", "", true)
		token := signTestJWT(t, secret, time.Now().Add(time.Hour))
		c.Request().Header.Set("Authorization", "Bearer "+token)

		var gotClaims any
		handler := mw(func(ctx *echo.Context) error {
			gotClaims = ctx.Request().Context().Value(jwtCtx.AuthClaimsKey)
			return ctx.NoContent(http.StatusOK)
		})

		if err := handler(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotClaims == nil {
			t.Fatal("expected auth claims in request context, got nil")
		}
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	t.Run("Return error when limiter.Take fails", func(t *testing.T) {
		mockLimiter := newLimiter(t)

		wantStatus := resPkg.NewStatusCode(http.StatusInternalServerError)
		mockLimiter.EXPECT().GetKeyREST(gomock.Any(), gomock.Any()).Return("rate-key")
		mockLimiter.EXPECT().Take(gomock.Any(), "rate-key").Return(nil, wantStatus)

		c, _ := newMiddlewareTestContext(t, http.MethodGet, "/api", "", true)
		handler := RateLimitMiddleware(mockLimiter)(func(ctx *echo.Context) error {
			t.Error("handler should not be called when Take returns error")
			return nil
		})

		err := handler(c)
		if err == nil {
			t.Fatal("expected error when Take fails, got nil")
		}
	})

	t.Run("Set headers and invoke next when limit is not reached", func(t *testing.T) {
		mockLimiter := newLimiter(t)

		mockLimiter.EXPECT().GetKeyREST(gomock.Any(), gomock.Any()).Return("rate-key")
		mockLimiter.EXPECT().Take(gomock.Any(), "rate-key").Return(&ratelimit.Result{
			Limit:     50,
			Remaining: 49,
			Reset:     30,
			Reached:   false,
		}, nil)

		c, rec := newMiddlewareTestContext(t, http.MethodGet, "/api", "", true)
		called := false
		handler := RateLimitMiddleware(mockLimiter)(func(ctx *echo.Context) error {
			called = true
			return ctx.NoContent(http.StatusOK)
		})

		if err := handler(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Error("expected next handler to be called")
		}
		if rec.Header().Get("X-RateLimit-Limit") != "50" {
			t.Errorf("X-RateLimit-Limit = %q, want %q", rec.Header().Get("X-RateLimit-Limit"), "50")
		}
		if rec.Header().Get("X-RateLimit-Remaining") != "49" {
			t.Errorf("X-RateLimit-Remaining = %q, want %q", rec.Header().Get("X-RateLimit-Remaining"), "49")
		}
		if rec.Header().Get("X-RateLimit-Reset") != "30" {
			t.Errorf("X-RateLimit-Reset = %q, want %q", rec.Header().Get("X-RateLimit-Reset"), "30")
		}
	})

	t.Run("Return rate-limit error with Retry-After when limit is reached", func(t *testing.T) {
		mockLimiter := newLimiter(t)

		mockLimiter.EXPECT().GetKeyREST(gomock.Any(), gomock.Any()).Return("rate-key")
		mockLimiter.EXPECT().Take(gomock.Any(), "rate-key").Return(&ratelimit.Result{
			Limit:      50,
			Remaining:  0,
			Reset:      30,
			RetryAfter: 15,
			Reached:    true,
		}, nil)

		c, rec := newMiddlewareTestContext(t, http.MethodGet, "/api", "", true)
		handler := RateLimitMiddleware(mockLimiter)(func(ctx *echo.Context) error {
			t.Error("handler should not be called when rate limit is reached")
			return nil
		})

		if err := handler(c); err == nil {
			t.Fatal("expected rate limit error, got nil")
		}
		if rec.Header().Get("Retry-After") != "15" {
			t.Errorf("Retry-After = %q, want %q", rec.Header().Get("Retry-After"), "15")
		}
	})

	t.Run("Return error when limit is reached and language context is missing", func(t *testing.T) {
		mockLimiter := newLimiter(t)

		mockLimiter.EXPECT().GetKeyREST(gomock.Any(), gomock.Any()).Return("rate-key")
		mockLimiter.EXPECT().Take(gomock.Any(), "rate-key").Return(&ratelimit.Result{
			Limit:      50,
			Remaining:  0,
			Reset:      30,
			RetryAfter: 15,
			Reached:    true,
		}, nil)

		c, _ := newMiddlewareTestContext(t, http.MethodGet, "/api", "", false)
		handler := RateLimitMiddleware(mockLimiter)(func(ctx *echo.Context) error {
			t.Error("handler should not be called when rate limit is reached")
			return nil
		})

		if err := handler(c); err == nil {
			t.Fatal("expected error when language context is missing, got nil")
		}
	})
}
