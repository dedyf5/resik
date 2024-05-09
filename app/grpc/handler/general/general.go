// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package general

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/request"
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
)

type GeneralHandler struct {
	log    *logCtx.Log
	config config.Config
}

func New(log *logCtx.Log, config config.Config) *GeneralHandler {
	return &GeneralHandler{
		log:    log,
		config: config,
	}
}

func (h *GeneralHandler) Home(c context.Context, common *request.Common) (*App, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err.GRPC().Err()
	}

	langReqCode := ""
	if ctx.Lang.LangReq != nil {
		langReqCode = ctx.Lang.LangReq.String()
	}

	lang := ctx.Lang

	return &App{
		App:     h.config.App.Name,
		Version: h.config.App.Version,
		Lang: &Lang{
			Current: lang.LanguageReqOrDefault().String(),
			Request: langReqCode,
			Default: h.config.App.LangDefault.String(),
		},
		Message: lang.GetByTemplateData("home_message", commonEntity.Map{"app_name": h.config.App.Name, "code": h.config.App.Version}),
	}, nil
}

func (h *GeneralHandler) mustEmbedUnimplementedGeneralServiceServer() {}
