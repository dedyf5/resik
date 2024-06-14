// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"context"
	"net/http"
	"strings"
	"time"

	jwtCtx "github.com/dedyf5/resik/ctx/jwt"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/entities/config"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"

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

func LoggerAndResponseFormatterMiddleware(log *logCtx.Log, appModule config.Module) echo.MiddlewareFunc {
	return echo.WrapMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if strings.Contains(r.RequestURI, DocPrefix) {
				h.ServeHTTP(w, r)
				return
			}

			correlationID := xid.New().String()
			ctx := context.WithValue(
				r.Context(),
				logCtx.CorrelationIDKeyContext,
				correlationID,
			)

			r = r.WithContext(ctx)

			// log.Logger = log.Logger.With(zap.String(string(correlationIDCtxKey), correlationID))
			log.CorrelationID = correlationID

			w.Header().Add(logCtx.CorrelationIDKeyXHeader.String(), correlationID)

			lrw := logCtx.NewHTTP(w, appModule, log, time.Now(), r.Method, r.RequestURI, r.UserAgent())

			r = r.WithContext(log.WithContext(ctx))

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
				return &resPkg.Status{
					Code:       http.StatusInternalServerError,
					CauseError: langErr,
				}
			}
			return jwtCtx.HTTPStatusError(err, langRes)
		},
	}
	return echojwt.WithConfig(jwtConfig)
}

func JWTMiddleware(signatureKey string) echo.MiddlewareFunc {
	return echo.WrapMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
			if token != "" {
				claim, _ := jwtCtx.AuthClaimsFromString(token, signatureKey)
				ctx := context.WithValue(r.Context(),
					jwtCtx.AuthClaimsKey,
					claim,
				)
				r = r.WithContext(ctx)
			}

			h.ServeHTTP(w, r)
		})
	})
}
