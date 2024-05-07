// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	generalHandler "github.com/dedyf5/resik/app/grpc/handler/general"
	"github.com/dedyf5/resik/config"
	"google.golang.org/grpc"
)

type Router struct {
	config         config.Config
	generalHandler *generalHandler.GeneralHandler
}

func newRouter(config config.Config, generalHandler *generalHandler.GeneralHandler) *Router {
	return &Router{
		config:         config,
		generalHandler: generalHandler,
	}
}

func (r *Router) routerSetup(grpcServer *grpc.Server) {
	generalHandler.RegisterGeneralServiceServer(grpcServer, r.generalHandler)
}
