// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"net/http"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/status"
	"github.com/labstack/echo/v4"
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
