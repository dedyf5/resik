// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package middleware

import (
	"net/http"

	statusEntity "github.com/dedyf5/resik/entities/status"
	httpUtil "github.com/dedyf5/resik/utils/http"
	langUtil "github.com/dedyf5/resik/utils/lang"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

func StatusEcho(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	switch status := err.(type) {
	case *statusEntity.HTTP:
		if status.Code != http.StatusNoContent {
			ctx.JSON(status.Code, httpUtil.ResponseFromStatusHTTP(status))
		}
		return
	}

	ctx.NoContent(http.StatusNoContent)
}

func LangEcho(langDefault language.Tag) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			langReq := langDefault
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
				langReq = *langRes
			}
			langAccept := c.Request().Header.Get("Accept-Language")
			c.Set(langKey, langUtil.NewLang(langDefault, langReq, langAccept))
			return next(c)
		}
	}
}
