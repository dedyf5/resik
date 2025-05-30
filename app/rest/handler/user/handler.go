// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"net/http"

	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	userService "github.com/dedyf5/resik/core/user"
	reqUserCore "github.com/dedyf5/resik/core/user/request"
	resUserCore "github.com/dedyf5/resik/core/user/response"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
	"github.com/dedyf5/resik/entities/user/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	fw          echoFW.IEcho
	log         *logCtx.Log
	userService userService.IService
}

func New(fw echoFW.IEcho, log *logCtx.Log, userService userService.IService) *Handler {
	return &Handler{
		log:         log,
		fw:          fw,
		userService: userService,
	}
}

// @Summary Login
// @Description Login by username and password
// @Tags user
// @Accept json
// @Produce json
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       payload body reqUserCore.LoginPost true "Payload"
// @Success		200	{object}	resPkg.Response{data=resUserCore.UserCredential}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     401 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/login [post]
func (h *Handler) LoginPost(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	ctx.App.Logger().Debug("LoginPost")

	var payload reqUserCore.LoginPost

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	token, err := h.userService.Auth(param.Auth{Ctx: ctx, Username: payload.Username, Password: payload.Password})
	if err != nil {
		return err
	}

	return &resPkg.Status{
		Code: http.StatusOK,
		Data: resUserCore.UserCredential{
			Username: payload.Username,
			Token:    token,
		},
	}
}

// @Summary Token Refresh
// @Description Request new token by existing token
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       parameter query commonEntity.Request true "Query Param"
// @Success		200	{object}	resPkg.Response{data=resUserCore.UserCredential}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/token-refresh [get]
func (h *Handler) TokenRefreshGet(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	ctx.App.Logger().Debug("TokenRefresh")

	var query commonEntity.Request
	if err := h.fw.StructValidator(echoCtx, &query); err != nil {
		return err
	}

	token, err := h.userService.AuthTokenGenerate(ctx.UserClaims.UserID, ctx.UserClaims.Username)
	if err != nil {
		return err
	}

	return &resPkg.Status{
		Code: http.StatusOK,
		Data: resUserCore.UserCredential{
			Username: ctx.UserClaims.Username,
			Token:    token,
		},
	}
}
