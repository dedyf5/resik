// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package handler

import (
	"net/http"

	"github.com/dedyf5/resik/config"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/status"
	commonEntity "github.com/dedyf5/resik/entities/common"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	config config.Config
}

func New(config config.Config) *Handler {
	return &Handler{
		config: config,
	}
}

func (h *Handler) Home(ctx echo.Context) error {
	lang := ctx.Get(langCtx.ContextKey.String()).(*langCtx.Lang)
	return &status.Status{
		Code:    http.StatusOK,
		Message: lang.GetByTemplateData("home_message", commonEntity.Map{"app_name": h.config.App.Name, "code": h.config.App.Version}),
		Data: commonEntity.Map{
			"app":          h.config.App.Name,
			"version":      h.config.App.Version,
			"current_lang": lang.LangReq,
			"default_lang": h.config.App.LangDefault,
		},
	}
}
