// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

//go:build wireinject
// +build wireinject

package bootstrap

import (
	generalHandler "github.com/dedyf5/resik/app/grpc/handler/general"
	trxHandler "github.com/dedyf5/resik/app/grpc/handler/transaction"
	"github.com/dedyf5/resik/app/grpc/middleware"
	"github.com/dedyf5/resik/config"
	trx "github.com/dedyf5/resik/core/transaction"
	trxService "github.com/dedyf5/resik/core/transaction/service"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/drivers"
	configEntity "github.com/dedyf5/resik/entities/config"
	repo "github.com/dedyf5/resik/repositories"
	trxRepo "github.com/dedyf5/resik/repositories/transaction"
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
	drivers.NewMySQLConnection,
	drivers.NewGorm,
)

var gormRepoSet = wire.NewSet(
	trxRepo.New,
	wire.Bind(new(repo.ITransaction), new(*trxRepo.TransactionRepo)),
)

var serviceSet = wire.NewSet(
	trxService.New,
	wire.Bind(new(trx.IService), new(*trxService.Service)),
)

var handlerSet = wire.NewSet(
	generalHandler.New,
	trxHandler.New,
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
		newServerHTTP,
		newRouter,
		newApp,
	)

	return nil, nil, nil
}
