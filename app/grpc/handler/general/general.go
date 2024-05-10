// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package general

import (
	"github.com/dedyf5/resik/config"
	logCtx "github.com/dedyf5/resik/ctx/log"
)

type GeneralHandler struct {
	config config.Config
	log    *logCtx.Log
}

func New(config config.Config, log *logCtx.Log) *GeneralHandler {
	return &GeneralHandler{
		log:    log,
		config: config,
	}
}

func (h *GeneralHandler) mustEmbedUnimplementedGeneralServiceServer() {}
