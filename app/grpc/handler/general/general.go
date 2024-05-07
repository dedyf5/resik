// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package general

import (
	"context"
	"log"

	"github.com/dedyf5/resik/config"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GeneralHandler struct {
	config config.Config
}

func New(config config.Config) *GeneralHandler {
	return &GeneralHandler{
		config: config,
	}
}

func (h *GeneralHandler) Home(ctx context.Context, _ *emptypb.Empty) (*App, error) {
	log.Print("Received request from client")
	return &App{
		App:     h.config.App.Name,
		Version: h.config.App.Version,
		Lang: &Lang{
			Current: h.config.App.LangDefault.String(),
			Request: "",
			Default: h.config.App.LangDefault.String(),
		},
	}, nil
}

func (h *GeneralHandler) mustEmbedUnimplementedGeneralServiceServer() {}
