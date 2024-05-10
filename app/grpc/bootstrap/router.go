// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	generalHandler "github.com/dedyf5/resik/app/grpc/handler/general"
	trxHandler "github.com/dedyf5/resik/app/grpc/handler/transaction"
	"github.com/dedyf5/resik/config"
	"google.golang.org/grpc"
)

type Router struct {
	config         config.Config
	generalHandler *generalHandler.GeneralHandler
	trxHandler     *trxHandler.TransactionHandler
}

func newRouter(config config.Config, generalHandler *generalHandler.GeneralHandler, trxHandler *trxHandler.TransactionHandler) *Router {
	return &Router{
		config:         config,
		generalHandler: generalHandler,
		trxHandler:     trxHandler,
	}
}

func (r *Router) routerSetup(grpcServer *grpc.Server) {
	generalHandler.RegisterGeneralServiceServer(grpcServer, r.generalHandler)
	trxHandler.RegisterTransactionServiceServer(grpcServer, r.trxHandler)
}
