// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package general

import (
	"github.com/dedyf5/resik/config"
	logCtx "github.com/dedyf5/resik/ctx/log"
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

func (h *GeneralHandler) mustEmbedUnimplementedGeneralServiceServer() {}
