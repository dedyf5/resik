// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

//go:build wireinject
// +build wireinject

package bootstrap

import (
	generalHandler "github.com/dedyf5/resik/app/grpc/handler/general"
	healthHandler "github.com/dedyf5/resik/app/grpc/handler/health"
	merchantHandler "github.com/dedyf5/resik/app/grpc/handler/merchant"
	trxHandler "github.com/dedyf5/resik/app/grpc/handler/transaction"
	userHandler "github.com/dedyf5/resik/app/grpc/handler/user"
	"github.com/dedyf5/resik/app/grpc/middleware"
	"github.com/dedyf5/resik/config"
	coreHealth "github.com/dedyf5/resik/core/health"
	dbChecker "github.com/dedyf5/resik/core/health/checkers"
	healthService "github.com/dedyf5/resik/core/health/service"
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

var configGeneral = config.Load(configEntity.ModuleGRPC)

var configGeneralSet = wire.NewSet(
	wire.Value(*configGeneral),
	wire.FieldsOf(new(config.Config), "APP", "HTTP", "Database", "Auth", "Log"),
	wire.FieldsOf(new(configEntity.App), "Env", "LangDefault"),
	wire.FieldsOf(new(drivers.SQLConfig), "Engine"),
)

var utilSet = wire.NewSet(
	validatorUtil.New,
	logCtx.Get,
)

var interceptorSet = wire.NewSet(
	middleware.NewInterceptor,
)

var connSet = wire.NewSet(
	wire.Value(false),
	drivers.NewMySQLConnection,
	drivers.NewGorm,
)

var gormRepoSet = wire.NewSet(
	merchantRepo.New,
	trxRepo.New,
	userRepo.New,
	wire.Bind(new(repo.IMerchant), new(*merchantRepo.MerchantRepo)),
	wire.Bind(new(repo.ITransaction), new(*trxRepo.TransactionRepo)),
	wire.Bind(new(repo.IUser), new(*userRepo.UserRepo)),
)

var serviceSet = wire.NewSet(
	merchantService.New,
	trxService.New,
	userService.New,
	healthService.New,
	wire.Bind(new(merchant.IService), new(*merchantService.Service)),
	wire.Bind(new(trx.IService), new(*trxService.Service)),
	wire.Bind(new(user.IService), new(*userService.Service)),
)

var handlerSet = wire.NewSet(
	generalHandler.New,
	merchantHandler.New,
	trxHandler.New,
	userHandler.New,
	healthHandler.New,
)

func provideCheckerSlice(dbChk coreHealth.Checker) []coreHealth.Checker {
	return []coreHealth.Checker{dbChk}
}

var healthCheckSet = wire.NewSet(
	dbChecker.NewDatabaseChecker,
	provideCheckerSlice,
)

func InitializeHTTP() (*App, func(), error) {
	wire.Build(
		configGeneralSet,
		utilSet,
		interceptorSet,
		connSet,
		gormRepoSet,
		serviceSet,
		handlerSet,
		healthCheckSet,
		newServerHTTP,
		newRouter,
		newApp,
	)

	return nil, nil, nil
}
