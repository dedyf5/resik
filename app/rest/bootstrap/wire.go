// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

//go:build wireinject
// +build wireinject

package bootstrap

import (
	"context"

	fw "github.com/dedyf5/resik/app/rest/fw/echo"
	generalHandler "github.com/dedyf5/resik/app/rest/handler/general"
	healthHandler "github.com/dedyf5/resik/app/rest/handler/health"
	merchantHandler "github.com/dedyf5/resik/app/rest/handler/merchant"
	trxHandler "github.com/dedyf5/resik/app/rest/handler/transaction"
	userHandler "github.com/dedyf5/resik/app/rest/handler/user"
	"github.com/dedyf5/resik/config"
	health "github.com/dedyf5/resik/core/health"
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
	pkgHash "github.com/dedyf5/resik/pkg/hash"
	repo "github.com/dedyf5/resik/repositories"
	checkRepo "github.com/dedyf5/resik/repositories/check"
	merchantRepo "github.com/dedyf5/resik/repositories/merchant"
	trxRepo "github.com/dedyf5/resik/repositories/transaction"
	userRepo "github.com/dedyf5/resik/repositories/user"
	ratelimitUtil "github.com/dedyf5/resik/utils/ratelimit"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
	"github.com/google/wire"
)

var configGeneral = config.Load(configEntity.ModuleTypeREST)

var configGeneralSet = wire.NewSet(
	wire.Value(*configGeneral),
	wire.FieldsOf(new(config.Config), "Module", "HTTP", "Database", "Redis", "RateLimit", "Auth", "Log"),
	wire.FieldsOf(new(configEntity.Module), "Env", "LangDefault", "Type"),
	wire.FieldsOf(new(drivers.SQLConfig), "Engine"),
)

var utilSet = wire.NewSet(
	validatorUtil.New,
	wire.Bind(new(validatorUtil.IValidate), new(*validatorUtil.Validate)),
	ratelimitUtil.NewRateLimiter,
	logCtx.Get,
)

var fwSet = wire.NewSet(
	fw.New,
	wire.Bind(new(fw.IEcho), new(*fw.Echo)),
)

var connSet = wire.NewSet(
	wire.Value(false),
	drivers.NewMySQLConnection,
	drivers.NewGorm,
	drivers.NewRedisConnection,
)

var repoSet = wire.NewSet(
	checkRepo.NewCheckDatabaseRepo,
	checkRepo.NewCheckRedisRepo,
	userRepo.New,
	merchantRepo.New,
	trxRepo.New,
	wire.Bind(new(repo.ICheck), new(*checkRepo.CheckDatabaseRepo)),
	wire.Bind(new(repo.IUser), new(*userRepo.UserRepo)),
	wire.Bind(new(repo.IMerchant), new(*merchantRepo.MerchantRepo)),
	wire.Bind(new(repo.ITransaction), new(*trxRepo.TransactionRepo)),
)

func provideHasherConfig(conf configEntity.Auth) *pkgHash.Argon2Config {
	return &pkgHash.Argon2Config{
		Memory:     conf.HashMemory,
		Iterations: conf.HashIterations,
	}
}

func provideCheckerSlice(db *checkRepo.CheckDatabaseRepo, redis *checkRepo.CheckRedisRepo) []repo.ICheck {
	return []repo.ICheck{db, redis}
}

var providerSet = wire.NewSet(
	provideHasherConfig,
	pkgHash.NewArgon2Hasher,
	provideCheckerSlice,
)

var serviceSet = wire.NewSet(
	userService.New,
	merchantService.New,
	trxService.New,
	healthService.New,
	wire.Bind(new(user.IService), new(*userService.Service)),
	wire.Bind(new(merchant.IService), new(*merchantService.Service)),
	wire.Bind(new(trx.IService), new(*trxService.Service)),
	wire.Bind(new(health.IService), new(*healthService.Service)),
)

var handlerSet = wire.NewSet(
	generalHandler.New,
	userHandler.New,
	merchantHandler.New,
	trxHandler.New,
	healthHandler.New,
)

func InitializeHTTP(c context.Context) (*App, func(), error) {
	wire.Build(
		configGeneralSet,
		utilSet,
		fwSet,
		connSet,
		repoSet,
		providerSet,
		serviceSet,
		handlerSet,
		newServerHTTP,
		newRouter,
		newApp,
	)

	return nil, nil, nil
}
