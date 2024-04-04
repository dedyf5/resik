// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

//go:build wireinject
// +build wireinject

package bootstrap

import (
	fw "github.com/dedyf5/resik/app/http/fw/echo"
	generalHandler "github.com/dedyf5/resik/app/http/handler/general"
	merchantHandler "github.com/dedyf5/resik/app/http/handler/merchant"
	trxHandler "github.com/dedyf5/resik/app/http/handler/transaction"
	userHandler "github.com/dedyf5/resik/app/http/handler/user"
	"github.com/dedyf5/resik/config"
	merchant "github.com/dedyf5/resik/core/merchant"
	merchantService "github.com/dedyf5/resik/core/merchant/service"
	trx "github.com/dedyf5/resik/core/transaction"
	trxService "github.com/dedyf5/resik/core/transaction/service"
	user "github.com/dedyf5/resik/core/user"
	userService "github.com/dedyf5/resik/core/user/service"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/drivers"
	configEntity "github.com/dedyf5/resik/entities/config"
	repo "github.com/dedyf5/resik/repositories"
	merchantRepo "github.com/dedyf5/resik/repositories/merchant"
	trxRepo "github.com/dedyf5/resik/repositories/transaction"
	userRepo "github.com/dedyf5/resik/repositories/user"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
	"github.com/google/wire"
)

var configGeneral = config.Load()

var configGeneralSet = wire.NewSet(
	wire.Value(*configGeneral),
	wire.FieldsOf(new(config.Config), "APP", "HTTP", "Database", "Log"),
	wire.FieldsOf(new(configEntity.App), "Env", "LangDefault"),
	wire.FieldsOf(new(drivers.SQLConfig), "Engine"),
)

var utilSet = wire.NewSet(
	validatorUtil.New,
	wire.Bind(new(validatorUtil.IValidate), new(*validatorUtil.Validate)),
	logCtx.Get,
)

var fwSet = wire.NewSet(
	fw.New,
	wire.Bind(new(fw.IEcho), new(*fw.Echo)),
)

var connSet = wire.NewSet(
	drivers.NewMySQLConnection,
	drivers.NewGorm,
)

var gormRepoSet = wire.NewSet(
	userRepo.New,
	merchantRepo.New,
	trxRepo.New,
	wire.Bind(new(repo.IUser), new(*userRepo.UserRepo)),
	wire.Bind(new(repo.IMerchant), new(*merchantRepo.MerchantRepo)),
	wire.Bind(new(repo.ITransaction), new(*trxRepo.TransactionRepo)),
)

var serviceSet = wire.NewSet(
	userService.New,
	merchantService.New,
	trxService.New,
	wire.Bind(new(user.IService), new(*userService.Service)),
	wire.Bind(new(merchant.IService), new(*merchantService.Service)),
	wire.Bind(new(trx.IService), new(*trxService.Service)),
)

var handlerSet = wire.NewSet(
	generalHandler.New,
	userHandler.New,
	merchantHandler.New,
	trxHandler.New,
)

func InitializeHTTP() (*App, func(), error) {
	wire.Build(
		configGeneralSet,
		utilSet,
		fwSet,
		connSet,
		gormRepoSet,
		serviceSet,
		handlerSet,
		newServerHTTP,
		newRouter,
		newApp,
	)

	return nil, nil, nil
}
