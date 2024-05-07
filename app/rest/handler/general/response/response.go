// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/ctx"
)

type Home struct {
	App     string `json:"app" example:"Resik"`
	Version string `json:"version" example:"0.1"`
	Lang    *Lang  `json:"lang"`
}

type Lang struct {
	Current string `json:"current" example:"id"`
	Request string `json:"request" example:"id"`
	Default string `json:"default" example:"en"`
}

func HomeMap(ctx *ctx.Ctx, config config.Config) Home {
	lang := ctx.Lang
	langReqCode := ""
	if lang.LangReq != nil {
		langReqCode = lang.LangReq.String()
	}
	return Home{
		App:     config.App.Name,
		Version: config.App.Version,
		Lang: &Lang{
			Current: lang.LanguageReqOrDefault().String(),
			Request: langReqCode,
			Default: config.App.LangDefault.String(),
		},
	}
}
