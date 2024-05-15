// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqUserCore "github.com/dedyf5/resik/core/user/request"
	resUserCore "github.com/dedyf5/resik/core/user/response"
	"github.com/dedyf5/resik/ctx"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *UserHandler) LoginPost(c context.Context, req *reqUserCore.LoginPost) (*UserCredentialRes, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err.GRPC().Err()
	}
	ctx.App.Logger().Debug("LoginPost")

	if err := h.validator.Struct(req, ctx.Lang); err != nil {
		return nil, err.GRPC().Err()
	}

	token, err := h.userService.Auth(req.ToParam(ctx))
	if err != nil {
		return nil, err.GRPC().Err()
	}

	return &UserCredentialRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: &resUserCore.UserCredential{
			Username: req.Username,
			Token:    token,
		},
	}, nil
}

func (h *UserHandler) TokenRefreshGet(c context.Context, _ *emptypb.Empty) (*UserCredentialRes, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err.GRPC().Err()
	}
	ctx.App.Logger().Debug("TokenRefreshGet")

	token, err := h.userService.AuthTokenGenerate(ctx.UserClaims.UserID, ctx.UserClaims.Username)
	if err != nil {
		return nil, err.GRPC().Err()
	}

	return &UserCredentialRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: &resUserCore.UserCredential{
			Username: ctx.UserClaims.Username,
			Token:    token,
		},
	}, nil
}
