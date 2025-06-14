// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	generalHandler "github.com/dedyf5/resik/app/grpc/handler/general"
	healthHandler "github.com/dedyf5/resik/app/grpc/handler/health"
	merchantHandler "github.com/dedyf5/resik/app/grpc/handler/merchant"
	trxHandler "github.com/dedyf5/resik/app/grpc/handler/transaction"
	userHandler "github.com/dedyf5/resik/app/grpc/handler/user"
	"github.com/dedyf5/resik/config"
	"google.golang.org/grpc"
)

type Router struct {
	config          config.Config
	generalHandler  *generalHandler.GeneralHandler
	merchantHandler *merchantHandler.MerchantHandler
	trxHandler      *trxHandler.TransactionHandler
	userHandler     *userHandler.UserHandler
	healthHandler   *healthHandler.HealthHandler
}

func newRouter(config config.Config, generalHandler *generalHandler.GeneralHandler, merchantHandler *merchantHandler.MerchantHandler, trxHandler *trxHandler.TransactionHandler, userHandler *userHandler.UserHandler, healthHandler *healthHandler.HealthHandler) *Router {
	return &Router{
		config:          config,
		generalHandler:  generalHandler,
		merchantHandler: merchantHandler,
		trxHandler:      trxHandler,
		userHandler:     userHandler,
		healthHandler:   healthHandler,
	}
}

func (r *Router) routerSetup(grpcServer *grpc.Server) {
	generalHandler.RegisterGeneralServiceServer(grpcServer, r.generalHandler)
	merchantHandler.RegisterMerchantServiceServer(grpcServer, r.merchantHandler)
	trxHandler.RegisterTransactionServiceServer(grpcServer, r.trxHandler)
	userHandler.RegisterUserServiceServer(grpcServer, r.userHandler)
	healthHandler.RegisterHealthServiceServer(grpcServer, r.healthHandler)
}
