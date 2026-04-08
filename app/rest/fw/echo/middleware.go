// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwtCtx "github.com/dedyf5/resik/ctx/jwt"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/ratelimit"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"golang.org/x/text/language"
)

func LangMiddleware(langDefault language.Tag) echo.MiddlewareFunc {
	return echo.WrapMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var langReq *language.Tag = nil
			langKey := langCtx.ContextKey.String()
			langString := r.URL.Query().Get(langKey)
			if langString != "" {
				langRes, err := langCtx.GetLanguageAvailable(langString)
				if err == nil {
					langReq = langRes
				}
			}
			ctx := context.WithValue(r.Context(),
				langCtx.ContextKey,
				langCtx.NewLang(langDefault,
					langReq,
					r.Header.Get("Accept-Language"),
				),
			)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	})
}

func LoggerAndResponseFormatterMiddleware(log *logCtx.Log) echo.MiddlewareFunc {
	return echo.WrapMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if strings.Contains(r.RequestURI, DocPrefix) {
				h.ServeHTTP(w, r)
				return
			}

			var requestBody []byte
			contentType := r.Header.Get("Content-Type")

			if !strings.HasPrefix(contentType, "multipart/form-data") && r.Body != nil && r.Body != http.NoBody {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					requestBody = bodyBytes
				}
			}

			correlationID, c := logCtx.EnsureCorrelationID(r.Context())

			callerHolder := &logCtx.CallerHolder{}

			c = context.WithValue(c, logCtx.KeyCallerHolderContext, callerHolder)

			r = r.WithContext(c)

			log.CorrelationID = correlationID
			log.Path = r.URL.Path
			log.QueryString = &r.URL.RawQuery

			lrw := logCtx.NewHTTP(w, c, log, time.Now(), r.Method, r.URL, contentType, r.UserAgent(), requestBody)

			h.ServeHTTP(lrw, r)
		})
	})
}

func ValidateTokenMiddleware(signatureKey string) echo.MiddlewareFunc {
	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCtx.AuthClaims)
		},
		SigningKey: []byte(signatureKey),
		ErrorHandler: func(c echo.Context, err error) error {
			ctx := c.Request().Context()
			langRes, langErr := langCtx.FromContext(ctx)
			if langErr != nil {
				return resPkg.NewStatusError(http.StatusInternalServerError, langErr)
			}
			return jwtCtx.HTTPStatusError(err, langRes)
		},
	}
	return echojwt.WithConfig(jwtConfig)
}

func JWTMiddleware(signatureKey string, langDef language.Tag) echo.MiddlewareFunc {
	return echo.WrapMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
			if token != "" {
				c := r.Context()
				langRes, langErr := langCtx.FromContext(c)
				if langErr != nil {
					langRes = langCtx.NewLang(langDef, nil, r.Header.Get("Accept-Language"))
				}
				claim, _ := jwtCtx.AuthClaimsFromString(token, signatureKey, langRes)
				ctx := context.WithValue(c,
					jwtCtx.AuthClaimsKey,
					claim,
				)
				r = r.WithContext(ctx)
			}

			h.ServeHTTP(w, r)
		})
	})
}

func RateLimitMiddleware(limiter ratelimit.Limiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := limiter.GetKeyREST(c.Request().Context(), c.RealIP())

			res, err := limiter.Take(c.Request().Context(), key)
			if err != nil {
				return err
			}

			w := c.Response().Header()
			w.Set("X-RateLimit-Limit", strconv.FormatInt(res.Limit, 10))
			w.Set("X-RateLimit-Remaining", strconv.FormatInt(res.Remaining, 10))
			w.Set("X-RateLimit-Reset", strconv.FormatInt(res.Reset, 10))

			if res.Reached {
				w.Set("Retry-After", strconv.FormatInt(res.RetryAfter, 10))

				lang, err := langCtx.FromContext(c.Request().Context())
				if err != nil {
					return err
				}

				return resPkg.NewStatusMessage(
					http.StatusTooManyRequests,
					lang.GetByMessageID("too_many_requests"),
					nil,
				)
			}

			return next(c)
		}
	}
}
