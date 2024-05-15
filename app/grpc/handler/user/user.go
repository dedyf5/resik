// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	userService "github.com/dedyf5/resik/core/user"
	logCtx "github.com/dedyf5/resik/ctx/log"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
)

type UserHandler struct {
	log         *logCtx.Log
	validator   *validatorUtil.Validate
	userService userService.IService
}

func New(log *logCtx.Log, validator *validatorUtil.Validate, userService userService.IService) *UserHandler {
	return &UserHandler{
		log:         log,
		validator:   validator,
		userService: userService,
	}
}

func (h *UserHandler) mustEmbedUnimplementedUserServiceServer() {}
