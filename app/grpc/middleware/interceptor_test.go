// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package middleware

import (
	"context"
	"net/http"
	"testing"
	"time"

	jwtCtx "github.com/dedyf5/resik/ctx/jwt"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	configEntity "github.com/dedyf5/resik/entities/config"
	mockIdentity "github.com/dedyf5/resik/internal/identity/mock"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/ratelimit"
	mockRatelimit "github.com/dedyf5/resik/utils/ratelimit/mock"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	loginPost          = "/user.UserService/LoginPost"
	orgListGet         = "/organization.OrganizationService/OrganizationListGet"
	testSignatureKey   = "test-signature-key"
	testTransportError = "no transport error"
)

type noopServerTransportStream struct {
	grpc.ServerTransportStream
}

func (noopServerTransportStream) SetHeader(metadata.MD) error {
	return nil
}

func (noopServerTransportStream) SendHeader(metadata.MD) error {
	return nil
}

func (noopServerTransportStream) SetTrailer(metadata.MD) error {
	return nil
}

func newTestInterceptorWithMocks(t *testing.T) (*Interceptor, *mockRatelimit.MockLimiter, *mockIdentity.MockIdentityResolver) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockLimiter := mockRatelimit.NewMockLimiter(ctrl)
	mockIdentity := mockIdentity.NewMockIdentityResolver(ctrl)
	logger := logCtx.Get(configEntity.Log{File: ""}, configEntity.Module{NameKey: "test_module"})

	ic := NewInterceptor(
		configEntity.Module{
			LangDefault: language.English,
		},
		configEntity.Auth{
			SignatureKey: "test-signature-key",
		},
		mockIdentity,
		mockLimiter,
		logger,
	)

	return ic, mockLimiter, mockIdentity
}

func withLanguageContext(c context.Context, langDefault language.Tag) context.Context {
	lang := langCtx.NewLang(langDefault, nil, "")
	return context.WithValue(c, langCtx.ContextKey, lang)
}

func signTestToken(t *testing.T, signatureKey string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub": "1",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(signatureKey))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return signed
}

func TestNewInterceptor_MethodRoles(t *testing.T) {
	ic, _, _ := newTestInterceptorWithMocks(t)
	if ic == nil {
		t.Fatal("expected non-nil Interceptor")
	}

	protectedMethods := []string{
		"/organization.OrganizationService/OrganizationDetailGet",
		orgListGet,
		"/transaction.TransactionService/TransactionPost",
		"/transaction.TransactionService/TransactionPut",
		"/transaction.TransactionService/TransactionDelete",
		"/transaction.TransactionService/TransactionDetailGet",
		"/transaction.TransactionService/TransactionListGet",
		"/transaction.TransactionService/OrganizationOmzetGet",
		"/transaction.TransactionService/BranchOmzetGet",
		"/user.UserService/TokenRefreshGet",
	}

	for _, method := range protectedMethods {
		roles, exists := ic.methodRoles[method]
		if !exists {
			t.Errorf("expected method %q to have role mapping", method)
			continue
		}
		if len(roles) != 1 || roles[0] != RoleValidToken {
			t.Errorf("method %q roles = %v, want [%v]", method, roles, RoleValidToken)
		}
	}
}

func TestInterceptor_RateLimit(t *testing.T) {
	ip := "ip:127.0.0.1"

	t.Run("Return error when language context is missing in RateLimit", func(t *testing.T) {
		ic, mockLimiter, _ := newTestInterceptorWithMocks(t)

		mockLimiter.EXPECT().GetKeyGRPC(gomock.Any()).Return(ip, nil)
		mockLimiter.EXPECT().
			Take(gomock.Any(), ip).
			Return(&ratelimit.Result{
				Limit:     100,
				Remaining: 99,
				Reset:     60,
				Reached:   false,
			}, nil)

		info := &grpc.UnaryServerInfo{FullMethod: loginPost}

		_, err := ic.RateLimit(context.Background(), "req", info, func(ctx context.Context, req any) (any, error) {
			t.Error("handler should not be called when lang context is missing")
			return nil, nil
		})
		if err == nil {
			t.Fatal("expected error when lang context is missing, got nil")
		}
	})

	t.Run("Bypass rate limit on HealthzGetPath", func(t *testing.T) {
		ic, _, _ := newTestInterceptorWithMocks(t)
		info := &grpc.UnaryServerInfo{FullMethod: HealthzGetPath}
		want := "health-ok"

		got, err := ic.RateLimit(context.Background(), "ping", info, func(ctx context.Context, req any) (any, error) {
			return want, nil
		})
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Allow request and set headers when rate limit not reached", func(t *testing.T) {
		ic, mockLimiter, _ := newTestInterceptorWithMocks(t)

		mockLimiter.EXPECT().GetKeyGRPC(gomock.Any()).Return(ip, nil)
		mockLimiter.EXPECT().
			Take(gomock.Any(), ip).
			Return(&ratelimit.Result{
				Limit:     100,
				Remaining: 99,
				Reset:     60,
				Reached:   false,
			}, nil)

		c := withLanguageContext(context.Background(), language.English)
		c = grpc.NewContextWithServerTransportStream(c, noopServerTransportStream{})
		info := &grpc.UnaryServerInfo{FullMethod: loginPost}
		want := "allowed-ok"

		got, err := ic.RateLimit(c, "req", info, func(ctx context.Context, req any) (any, error) {
			return want, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Return error when GetKeyGRPC fails", func(t *testing.T) {
		ic, mockLimiter, _ := newTestInterceptorWithMocks(t)

		wantStatus := resPkg.NewStatusCode(http.StatusInternalServerError)
		mockLimiter.EXPECT().GetKeyGRPC(gomock.Any()).Return("", wantStatus)

		info := &grpc.UnaryServerInfo{FullMethod: loginPost}

		_, err := ic.RateLimit(context.Background(), nil, info, func(ctx context.Context, req any) (any, error) {
			t.Error("handler should not be called")
			return nil, nil
		})
		if err == nil {
			t.Fatal("expected error when GetKeyGRPC fails, got nil")
		}
	})

	t.Run("Return error when Take fails", func(t *testing.T) {
		ic, mockLimiter, _ := newTestInterceptorWithMocks(t)

		mockLimiter.EXPECT().GetKeyGRPC(gomock.Any()).Return(ip, nil)

		wantStatus := resPkg.NewStatusCode(http.StatusInternalServerError)
		mockLimiter.EXPECT().Take(gomock.Any(), ip).Return(nil, wantStatus)

		info := &grpc.UnaryServerInfo{FullMethod: loginPost}

		_, err := ic.RateLimit(context.Background(), nil, info, func(ctx context.Context, req any) (any, error) {
			t.Error("handler should not be called")
			return nil, nil
		})
		if err == nil {
			t.Fatal("expected error when Take fails, got nil")
		}
	})

	t.Run("Reject request when rate limit reached", func(t *testing.T) {
		ic, mockLimiter, _ := newTestInterceptorWithMocks(t)

		mockLimiter.EXPECT().GetKeyGRPC(gomock.Any()).Return(ip, nil)
		mockLimiter.EXPECT().
			Take(gomock.Any(), ip).
			Return(&ratelimit.Result{
				Limit:      10,
				Remaining:  0,
				Reset:      30,
				RetryAfter: 15,
				Reached:    true,
			}, nil)

		c := withLanguageContext(context.Background(), language.Indonesian)
		info := &grpc.UnaryServerInfo{FullMethod: loginPost}

		_, err := ic.RateLimit(c, nil, info, func(c context.Context, req any) (any, error) {
			t.Error("handler should not be called when rate limit is reached")
			return nil, nil
		})
		if err == nil {
			t.Fatal("expected error when rate limit reached, got nil")
		}
	})

	t.Run("Return error when SetHeader fails without server transport stream", func(t *testing.T) {
		ic, mockLimiter, _ := newTestInterceptorWithMocks(t)

		mockLimiter.EXPECT().GetKeyGRPC(gomock.Any()).Return(ip, nil)
		mockLimiter.EXPECT().
			Take(gomock.Any(), ip).
			Return(&ratelimit.Result{
				Limit:     100,
				Remaining: 99,
				Reset:     60,
				Reached:   false,
			}, nil)

		c := withLanguageContext(context.Background(), language.English)
		info := &grpc.UnaryServerInfo{FullMethod: loginPost}
		_, err := ic.RateLimit(c, nil, info, func(ctx context.Context, req any) (any, error) {
			t.Error("handler should not be called when SetHeader fails")
			return nil, nil
		})
		if err == nil {
			t.Fatal("expected error when SetHeader fails without server transport stream, got nil")
		}
	})
}

func TestInterceptor_ValidateJWT_MissingMetadata(t *testing.T) {
	ic, _, _ := newTestInterceptorWithMocks(t)
	lang := langCtx.NewLang(language.English, &language.Indonesian, "id")

	_, err := ic.validateJWT(context.Background(), lang, testSignatureKey, loginPost)
	if err == nil {
		t.Fatal("expected error when metadata is missing in validateJWT, got nil")
	}
}

func TestInterceptor_Unary(t *testing.T) {
	t.Run("Fail when incoming metadata is missing", func(t *testing.T) {
		ic, _, _ := newTestInterceptorWithMocks(t)
		info := &grpc.UnaryServerInfo{FullMethod: loginPost}

		_, err := ic.Unary(context.Background(), "req", info, func(ctx context.Context, req any) (any, error) {
			t.Error("handler should not be called without metadata")
			return nil, nil
		})
		if err == nil {
			t.Fatal("expected error when metadata is missing, got nil")
		}
	})

	t.Run("Fail when language metadata is invalid", func(t *testing.T) {
		ic, _, _ := newTestInterceptorWithMocks(t)
		md := metadata.Pairs(langCtx.ContextKey.String(), "unsupported-lang-xyz")
		c := metadata.NewIncomingContext(context.Background(), md)
		info := &grpc.UnaryServerInfo{FullMethod: loginPost}

		_, err := ic.Unary(c, "req", info, func(ctx context.Context, req any) (any, error) {
			t.Error("handler should not be called when language is invalid")
			return nil, nil
		})
		if err == nil {
			t.Error("expected error for invalid language metadata, got nil")
		}
	})

	t.Run("Fail with unauthorized when RoleValidToken method has no bearer token", func(t *testing.T) {
		ic, _, _ := newTestInterceptorWithMocks(t)
		c := metadata.NewIncomingContext(context.Background(), metadata.Pairs())
		info := &grpc.UnaryServerInfo{FullMethod: orgListGet}

		_, err := ic.Unary(c, "req", info, func(c context.Context, req any) (any, error) {
			t.Error("handler should not be called without authorization header")
			return nil, nil
		})
		if err == nil {
			t.Fatal("expected unauthorized error, got nil")
		}
	})

	t.Run("Invoke handler when method roles does not contain RoleValidToken", func(t *testing.T) {
		ic, _, _ := newTestInterceptorWithMocks(t)

		customMethod := "/custom.CustomService/MethodWithoutValidTokenRole"
		ic.methodRoles[customMethod] = []Role{}

		c := metadata.NewIncomingContext(context.Background(), metadata.Pairs())
		info := &grpc.UnaryServerInfo{FullMethod: customMethod}

		wantRes := "custom-method-ok"
		gotRes, err := ic.Unary(c, "req", info, func(c context.Context, req any) (any, error) {
			return wantRes, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotRes != wantRes {
			t.Fatalf("gotRes = %v, want %v", gotRes, wantRes)
		}
	})

	t.Run("Fail when bearer token is invalid", func(t *testing.T) {
		ic, _, _ := newTestInterceptorWithMocks(t)
		c := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer invalid-token"))
		info := &grpc.UnaryServerInfo{FullMethod: orgListGet}

		_, err := ic.Unary(c, "req", info, func(ctx context.Context, req any) (any, error) {
			t.Error("handler should not be called with invalid token")
			return nil, nil
		})
		if err == nil {
			t.Fatal("expected error for invalid bearer token, got nil")
		}
	})

	t.Run("Invoke handler when bearer token is valid", func(t *testing.T) {
		ic, _, mockIdentity := newTestInterceptorWithMocks(t)

		mockIdentity.EXPECT().
			GetTenantOrganizationIDs(gomock.Any(), gomock.Any()).
			Return([]uint64{1}, nil)
		mockIdentity.EXPECT().
			GetTenantBranchIDs(gomock.Any(), gomock.Any()).
			Return([]uint64{1}, nil)

		token := signTestToken(t, testSignatureKey)
		md := metadata.Pairs(
			langCtx.ContextKey.String(), "en",
			"authorization", "Bearer "+token,
		)
		c := metadata.NewIncomingContext(context.Background(), md)
		info := &grpc.UnaryServerInfo{FullMethod: orgListGet}

		wantRes := "org-list-ok"
		gotRes, err := ic.Unary(c, "req", info, func(c context.Context, req any) (any, error) {
			if c.Value(jwtCtx.AuthClaimsKey) == nil {
				t.Error("expected AuthClaimsKey to be set in context")
			}
			return wantRes, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotRes != wantRes {
			t.Errorf("gotRes = %v, want %v", gotRes, wantRes)
		}
	})

	t.Run("Invoke handler for unprotected method with valid metadata", func(t *testing.T) {
		ic, _, _ := newTestInterceptorWithMocks(t)
		c := metadata.NewIncomingContext(context.Background(), metadata.Pairs())
		info := &grpc.UnaryServerInfo{FullMethod: loginPost}

		wantRes := "login-success"
		gotRes, err := ic.Unary(c, "login-req", info, func(c context.Context, req any) (any, error) {
			if ic.log.Path != info.FullMethod {
				t.Errorf("log.Path = %q, want %q", ic.log.Path, info.FullMethod)
			}
			if ic.log.CorrelationID == "" {
				t.Error("expected CorrelationID to be populated")
			}
			return wantRes, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotRes != wantRes {
			t.Errorf("gotRes = %v, want %v", gotRes, wantRes)
		}
	})
}
