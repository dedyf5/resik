// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	"github.com/dedyf5/resik/buildinfo"
	"github.com/dedyf5/resik/config"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/pkg/color"
	"github.com/labstack/echo/v5"
	echoMiddle "github.com/labstack/echo/v5/middleware"
)

type ServerHTTP struct {
	config     config.Config
	echo       *echo.Echo
	httpServer *http.Server
}

func newServerHTTP(config config.Config, log *logCtx.Log) *ServerHTTP {
	e := echo.New()
	e.Binder = echoFW.NewBinder()
	e.HTTPErrorHandler = echoFW.HTTPErrorHandler
	e.IPExtractor = func(r *http.Request) string {
		if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
			return ip
		}
		return echo.ExtractIPDirect()(r)
	}
	e.Use(echoFW.LoggerAndResponseFormatterMiddleware(log))
	e.Use(echoFW.LangMiddleware(config.Module.LangDefault))
	e.Use(echoMiddle.CORSWithConfig(
		echoMiddle.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		},
	))

	return &ServerHTTP{
		echo:   e,
		config: config,
	}
}

func (s *ServerHTTP) Start(c context.Context) {
	appName := color.Format(color.GREEN, s.config.App.Name())
	appVersion := color.Format(color.YELLOW, s.config.App.Version())

	addr := s.config.Module.HostPort()

	fmt.Printf("%s\n\n", buildinfo.FrameworkLogoASCIIVersion)
	fmt.Printf("%s version %s\n", appName, appVersion)
	fmt.Printf("Module %s type %s\n\n", s.config.Module.Name, s.config.Module.Type.String())
	log.Printf("STARTED HTTP SERVER AT %v\n", addr)

	go func() {
		s.httpServer = &http.Server{
			Addr:              addr,
			ReadHeaderTimeout: s.config.HTTP.ReadHeaderTimeout,
			ReadTimeout:       s.config.HTTP.ReadTimeout,
			WriteTimeout:      s.config.HTTP.WriteTimeout,
			IdleTimeout:       s.config.HTTP.IdleTimeout,
			BaseContext: func(_ net.Listener) context.Context {
				return c
			},
			Handler: s.echo,
		}

		err := s.httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("HTTP SERVER CLOSED")
			} else {
				log.Printf("HTTP SERVER ERROR: %s", err.Error())
			}
		}
	}()
}

func (s *ServerHTTP) Close() error {
	ctxShutdown, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	log.Println("SHUTTING DOWN HTTP SERVER")
	err := s.httpServer.Shutdown(ctxShutdown)
	if err != nil {
		log.Printf("HTTP SERVER ERROR: %s", err.Error())
	}
	log.Println("SUCCESSFULLY SHUTDOWN HTTP SERVER")
	return nil
}
