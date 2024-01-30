// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"context"
	"net/http"
	"time"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/ctx/status"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"

	"golang.org/x/text/language"
)

func LangMiddleware(langDefault language.Tag) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var langReq *language.Tag = nil
			langKey := langCtx.ContextKey.String()
			langString := c.Request().URL.Query().Get(langKey)
			if langString != "" {
				langRes, err := langCtx.GetLanguageAvailable(langString)
				if err != nil {
					return &status.Status{
						Code:    http.StatusBadRequest,
						Message: err.Error(),
						Detail: map[string]string{
							langKey: err.Error(),
						},
					}
				}
				langReq = langRes
			}
			langAccept := c.Request().Header.Get("Accept-Language")
			c.Set(langKey, langCtx.NewLang(langDefault, langReq, langAccept))
			return next(c)
		}
	}
}

func LoggerMiddleware(log *logCtx.Log) echo.MiddlewareFunc {
	return echo.WrapMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			correlationID := xid.New().String()
			ctx := context.WithValue(
				r.Context(),
				logCtx.CorrelationIDKeyContext,
				correlationID,
			)

			r = r.WithContext(ctx)

			// log.Logger = log.Logger.With(zap.String(string(correlationIDCtxKey), correlationID))
			log.CorrelationID = correlationID

			w.Header().Add("X-Correlation-ID", correlationID)

			lrw := logCtx.NewHTTP(w, log, time.Now(), r.Method, r.RequestURI, r.UserAgent())

			r = r.WithContext(log.WithContext(ctx))

			h.ServeHTTP(lrw, r)
		})
	})
}
