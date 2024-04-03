// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	"fmt"

	"github.com/dedyf5/resik/app/http/docs"
	echoFW "github.com/dedyf5/resik/app/http/fw/echo"
	generalHandler "github.com/dedyf5/resik/app/http/handler/general"
	trxHandler "github.com/dedyf5/resik/app/http/handler/transaction"
	userHandler "github.com/dedyf5/resik/app/http/handler/user"
	"github.com/dedyf5/resik/config"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	config         config.Config
	generalHandler *generalHandler.Handler
	userHandler    *userHandler.Handler
	trxHandler     *trxHandler.Handler
}

func newRouter(config config.Config, generalHandler *generalHandler.Handler, userHandler *userHandler.Handler, trxHandler *trxHandler.Handler) *Router {
	return &Router{
		config:         config,
		generalHandler: generalHandler,
		userHandler:    userHandler,
		trxHandler:     trxHandler,
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
	e.POST("/token-refresh", userHandler.TokenRefresh, validateToken, jwtMiddleware)

	trxHandler := r.trxHandler
	trx := e.Group("/transaction")
	trxMerchant := trx.Group("/merchant/:id", validateToken, jwtMiddleware)
	trxMerchant.GET("/omzet", trxHandler.MerchantOmzetGet)
	trxOutlet := trx.Group("/outlet/:id", validateToken, jwtMiddleware)
	trxOutlet.GET("/omzet", trxHandler.OutletOmzetGet)

	docs.SwaggerInfohttp.Title = r.config.App.Name
	docs.SwaggerInfohttp.Version = r.config.App.Version
	docs.SwaggerInfohttp.Host = r.config.App.HostPort()
	docs.SwaggerInfohttp.Description = r.config.App.APIDocDescription()
	docHandler := echoSwagger.EchoWrapHandler(
		echoSwagger.InstanceName(docs.SwaggerInfohttp.InfoInstanceName),
	)
	e.GET(fmt.Sprintf("%s*", echoFW.DocPrefix), docHandler)
}
