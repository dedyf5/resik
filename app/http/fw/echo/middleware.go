// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"net/http"

	statusEntity "github.com/dedyf5/resik/entities/status"
	langUtil "github.com/dedyf5/resik/utils/lang"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

func LangMiddleware(langDefault language.Tag) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var langReq *language.Tag = nil
			langKey := langUtil.ContextKey.String()
			langString := c.Request().URL.Query().Get(langKey)
			if langString != "" {
				langRes, err := langUtil.GetLanguageAvailable(langString)
				if err != nil {
					return &statusEntity.HTTP{
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
			c.Set(langKey, langUtil.NewLang(langDefault, langReq, langAccept))
			return next(c)
		}
	}
}
