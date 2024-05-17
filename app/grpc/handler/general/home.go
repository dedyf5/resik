// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package general

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	resAppCore "github.com/dedyf5/resik/core/app/response"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/common"
	"google.golang.org/grpc/codes"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (h *GeneralHandler) Home(c context.Context, _ *emptypb.Empty) (*HomeRes, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err
	}

	return &HomeRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: ctx.Lang.GetByTemplateData("home_message", common.Map{"app_name": h.config.App.Name, "code": h.config.App.Version}),
		},
		Data: resAppCore.AppMap(ctx, h.config),
	}, nil
}
