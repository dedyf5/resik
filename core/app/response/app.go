// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/common"
)

func AppMap(ctx *ctx.Ctx, config config.Config, req *common.Request) *App {
	lang := ctx.Lang
	langReqCode := ""
	if lang.LangReq != nil {
		langReqCode = lang.LangReq.String()
	}
	return &App{
		App:     config.App.Name,
		Version: config.App.Version,
		Lang: &AppLang{
			Current:   lang.LanguageReqOrDefault().String(),
			Request:   langReqCode,
			Default:   config.App.LangDefault.String(),
			Available: req.LangAvailable(),
		},
	}
}
