// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

//go:build wireinject
// +build wireinject

package bootstrap

import (
	generalHandler "github.com/dedyf5/resik/app/grpc/handler/general"
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/drivers"
	configEntity "github.com/dedyf5/resik/entities/config"
	"github.com/google/wire"
)

var configGeneral = config.Load(configEntity.ModuleGRPC)

var configGeneralSet = wire.NewSet(
	wire.Value(*configGeneral),
	wire.FieldsOf(new(config.Config), "APP", "HTTP", "Database", "Log"),
	wire.FieldsOf(new(configEntity.App), "Env", "LangDefault"),
	wire.FieldsOf(new(drivers.SQLConfig), "Engine"),
)

var handlerSet = wire.NewSet(
	generalHandler.New,
)

func InitializeHTTP() (*App, func(), error) {
	wire.Build(
		configGeneralSet,
		handlerSet,
		newServerHTTP,
		newRouter,
		newApp,
	)

	return nil, nil, nil
}
