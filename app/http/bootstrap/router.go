// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	generalHandler "github.com/dedyf5/resik/app/http/handler/general"
	trxHandler "github.com/dedyf5/resik/app/http/handler/transaction"
	"github.com/dedyf5/resik/config"
)

type Router struct {
	config         config.Config
	generalHandler *generalHandler.Handler
	trxHandler     *trxHandler.Handler
}

func newRouter(config config.Config, generalHandler *generalHandler.Handler, trxHandler *trxHandler.Handler) *Router {
	return &Router{
		config:         config,
		generalHandler: generalHandler,
		trxHandler:     trxHandler,
	}
}

func (r *Router) routerSetup(server *ServerHTTP) {
	e := server.echo

	generalHandler := r.generalHandler
	e.GET("/", generalHandler.Home)

	trxHandler := r.trxHandler
	trx := e.Group("/transaction")
	trx.POST("", trxHandler.Create)
	trxMerchant := trx.Group("/merchant/:id")
	trxMerchant.GET("/omzet", trxHandler.MerchantOmzetGet)
	trxOutlet := trx.Group("/outlet/:id")
	trxOutlet.GET("/omzet", trxHandler.OutletOmzetGet)
}
