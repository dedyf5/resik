// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package handler

import (
	"net/http"

	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	log    *logCtx.Log
	config config.Config
}

func New(log *logCtx.Log, config config.Config) *Handler {
	return &Handler{
		log:    log,
		config: config,
	}
}

func (h *Handler) Home(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	lang := ctx.Lang
	langReqCode := ""
	if lang.LangReq != nil {
		langReqCode = lang.LangReq.String()
	}
	return &resPkg.Status{
		Code:    http.StatusOK,
		Message: lang.GetByTemplateData("home_message", commonEntity.Map{"app_name": h.config.App.Name, "code": h.config.App.Version}),
		Data: commonEntity.Map{
			"app":     h.config.App.Name,
			"version": h.config.App.Version,
			"lang": commonEntity.Map{
				"current": lang.LanguageReqOrDefault().String(),
				"request": langReqCode,
				"default": h.config.App.LangDefault,
			},
		},
	}
}
