// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

//go:build wireinject
// +build wireinject

package bootstrap

import (
	generalHandler "github.com/dedyf5/resik/app/http/handler/general"
	trxHandler "github.com/dedyf5/resik/app/http/handler/transaction"
	"github.com/dedyf5/resik/config"
	trxService "github.com/dedyf5/resik/core/transaction/service"
	"github.com/dedyf5/resik/drivers"
	configEntity "github.com/dedyf5/resik/entities/config"
	trxRepo "github.com/dedyf5/resik/repositories/transaction"
	"github.com/google/wire"
)

var configGeneral = config.Load()

var configGeneralSet = wire.NewSet(
	wire.Value(*configGeneral),
	wire.FieldsOf(new(config.Config), "APP", "HTTP", "Database"),
	wire.FieldsOf(new(configEntity.App), "Env", "LangDefault"),
	wire.FieldsOf(new(drivers.SQLConfig), "Engine"),
)

var connSet = wire.NewSet(
	drivers.NewMySQLConnection,
	drivers.NewGorm,
)

var gormRepoSet = wire.NewSet(
	trxRepo.New,
	wire.Bind(new(trxRepo.ITransaction), new(*trxRepo.TransactionRepo)),
)

var serviceSet = wire.NewSet(
	trxService.New,
	wire.Bind(new(trxService.IService), new(*trxService.Service)),
)

var handlerSet = wire.NewSet(
	generalHandler.New,
	trxHandler.New,
)

func InitializeHTTP() (*App, func(), error) {
	wire.Build(
		configGeneralSet,
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
