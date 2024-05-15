// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	"fmt"

	"github.com/dedyf5/resik/app/rest/docs"
	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	generalHandler "github.com/dedyf5/resik/app/rest/handler/general"
	merchantHandler "github.com/dedyf5/resik/app/rest/handler/merchant"
	trxHandler "github.com/dedyf5/resik/app/rest/handler/transaction"
	userHandler "github.com/dedyf5/resik/app/rest/handler/user"
	"github.com/dedyf5/resik/config"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	config          config.Config
	generalHandler  *generalHandler.Handler
	merchantHandler *merchantHandler.Handler
	userHandler     *userHandler.Handler
	trxHandler      *trxHandler.Handler
}

func newRouter(config config.Config, generalHandler *generalHandler.Handler, userHandler *userHandler.Handler, merchantHandler *merchantHandler.Handler, trxHandler *trxHandler.Handler) *Router {
	return &Router{
		config:          config,
		generalHandler:  generalHandler,
		userHandler:     userHandler,
		merchantHandler: merchantHandler,
		trxHandler:      trxHandler,
	}
}

func (r *Router) routerSetup(server *ServerHTTP) {
	e := server.echo

	validateToken := echoFW.ValidateTokenMiddleware(r.config.Auth.SignatureKey)
	jwtMiddleware := echoFW.JWTMiddleware(r.config.Auth.SignatureKey)

	generalHandler := r.generalHandler
	e.GET("/", generalHandler.Home)

	userHandler := r.userHandler
	e.POST("/login", userHandler.LoginPost)
	e.GET("/token-refresh", userHandler.TokenRefreshGet, validateToken, jwtMiddleware)

	merchantHandler := r.merchantHandler
	e.GET("/merchant", merchantHandler.MerchantListGet, validateToken, jwtMiddleware)
	e.POST("/merchant", merchantHandler.MerchantPost, validateToken, jwtMiddleware)
	e.PUT("/merchant/:id", merchantHandler.MerchantPut, validateToken, jwtMiddleware)
	e.DELETE("/merchant/:id", merchantHandler.MerchantDelete, validateToken, jwtMiddleware)

	trxHandler := r.trxHandler
	trx := e.Group("/transaction")
	trxMerchant := trx.Group("/merchant/:id", validateToken, jwtMiddleware)
	trxMerchant.GET("/omzet", trxHandler.MerchantOmzetGet)
	trxOutlet := trx.Group("/outlet/:id", validateToken, jwtMiddleware)
	trxOutlet.GET("/omzet", trxHandler.OutletOmzetGet)

	docs.SwaggerInforest.Title = r.config.App.Name
	docs.SwaggerInforest.Version = r.config.App.Version
	docs.SwaggerInforest.Host = r.config.App.HostPort()
	docs.SwaggerInforest.Description = r.config.App.APIDocDescription()
	docHandler := echoSwagger.EchoWrapHandler(
		echoSwagger.InstanceName(docs.SwaggerInforest.InfoInstanceName),
	)
	e.GET(fmt.Sprintf("%s*", echoFW.DocPrefix), docHandler)
}
