// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package general

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	resAppCore "github.com/dedyf5/resik/core/app/response"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	"github.com/dedyf5/resik/entities/common"
	"google.golang.org/grpc/codes"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (h *GeneralHandler) Home(c context.Context, _ *emptypb.Empty) (*HomeRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}

	req := common.Request{}

	return &HomeRes{
		Status: &status.Status{
			Code: status.CodePlus(codes.OK),
			Message: term.HomeMessage.Localize(
				ctx.Lang().Localizer,
				h.config.App.Name(),
				h.config.App.Version(),
				h.config.Module.Name,
				h.config.Module.Type.String(),
			),
		},
		Data: resAppCore.AppMap(ctx, h.config, &req),
	}, nil
}
