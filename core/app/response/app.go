// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"github.com/dedyf5/resik/build"
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/common"
)

func AppMap(ctx *ctx.Ctx, config config.Config, req *common.Request) *App {
	lang := ctx.Lang()
	langReqCode := ""
	if lang.LangReq != nil {
		langReqCode = lang.LangReq.String()
	}
	return &App{
		App:     build.AppName,
		Version: build.AppVersion,
		Module: &Module{
			Name:    config.Module.Name,
			Type:    config.Module.Type.String(),
			Version: config.Module.Version,
		},
		Lang: &AppLang{
			Current:   lang.LanguageReqOrDefault().String(),
			Request:   langReqCode,
			Default:   config.Module.LangDefault.String(),
			Available: req.LangAvailable(),
		},
	}
}
