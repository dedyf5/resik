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
