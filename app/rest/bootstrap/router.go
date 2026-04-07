// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	"github.com/dedyf5/resik/app/rest/docs"
	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	generalHandler "github.com/dedyf5/resik/app/rest/handler/general"
	healthHandler "github.com/dedyf5/resik/app/rest/handler/health"
	merchantHandler "github.com/dedyf5/resik/app/rest/handler/merchant"
	trxHandler "github.com/dedyf5/resik/app/rest/handler/transaction"
	userHandler "github.com/dedyf5/resik/app/rest/handler/user"
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/utils/ratelimit"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	config          config.Config
	limiter         ratelimit.Limiter
	generalHandler  *generalHandler.Handler
	merchantHandler *merchantHandler.Handler
	userHandler     *userHandler.Handler
	trxHandler      *trxHandler.Handler
	healthHandler   *healthHandler.HealthHandler
}

func newRouter(config config.Config, limiter ratelimit.Limiter, generalHandler *generalHandler.Handler, userHandler *userHandler.Handler, merchantHandler *merchantHandler.Handler, trxHandler *trxHandler.Handler, healthHandler *healthHandler.HealthHandler) *Router {
	return &Router{
		config:          config,
		limiter:         limiter,
		generalHandler:  generalHandler,
		userHandler:     userHandler,
		merchantHandler: merchantHandler,
		trxHandler:      trxHandler,
		healthHandler:   healthHandler,
	}
}

func (r *Router) routerSetup(server *ServerHTTP) {
	e := server.echo

	validateToken := echoFW.ValidateTokenMiddleware(r.config.Auth.SignatureKey)
	jwtMiddleware := echoFW.JWTMiddleware(r.config.Auth.SignatureKey, r.config.Module.LangDefault)
	rateLimit := echoFW.RateLimitMiddleware(r.limiter)

	generalHandler := r.generalHandler
	e.GET("/", generalHandler.Home, rateLimit)

	userHandler := r.userHandler
	e.POST("/login", userHandler.LoginPost, rateLimit)
	e.GET("/token-refresh", userHandler.TokenRefreshGet, validateToken, jwtMiddleware, rateLimit)

	merchantHandler := r.merchantHandler
	merchant := e.Group("/merchant", validateToken, jwtMiddleware, rateLimit)
	merchant.GET("", merchantHandler.MerchantListGet)
	merchant.POST("", merchantHandler.MerchantPost)
	merchant.PUT("/:id", merchantHandler.MerchantPut)
	merchant.DELETE("/:id", merchantHandler.MerchantDelete)

	trxHandler := r.trxHandler
	trx := e.Group("/transaction", validateToken, jwtMiddleware, rateLimit)
	trxMerchant := trx.Group("/merchant/:merchant_id")
	trxMerchant.GET("/omzet", trxHandler.MerchantOmzetGet)
	trxOutlet := trx.Group("/outlet/:outlet_id")
	trxOutlet.GET("/omzet", trxHandler.OutletOmzetGet)

	healthH := r.healthHandler
	e.GET("/healthz", healthH.HealthHealthzGet)
	e.GET("/readyz", healthH.HealthReadyzGet, rateLimit)

	docs.SwaggerInforest.Title = r.config.Module.Name
	docs.SwaggerInforest.Version = r.config.Module.Version
	docs.SwaggerInforest.Host = r.config.Module.Public.HostPort()
	docs.SwaggerInforest.Schemes = []string{r.config.Module.Public.Schema}
	docs.SwaggerInforest.BasePath = r.config.Module.Public.BasePath
	docs.SwaggerInforest.Description = r.config.Module.APIDocDescription()
	docHandler := echoSwagger.EchoWrapHandler(
		echoSwagger.InstanceName(docs.SwaggerInforest.InfoInstanceName),
	)
	e.GET(echoFW.DocPrefix+"*", docHandler)
}
