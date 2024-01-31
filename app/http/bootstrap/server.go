// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	echoFW "github.com/dedyf5/resik/app/http/fw/echo"
	"github.com/dedyf5/resik/cmd"
	"github.com/dedyf5/resik/config"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/pkg/color"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

type ServerHTTP struct {
	config config.Config
	echo   *echo.Echo
}

func newServerHTTP(config config.Config, lang language.Tag, log *logCtx.Log) *ServerHTTP {
	echo := echo.New()
	echo.HideBanner = true
	echo.HidePort = true
	echo.Binder = echoFW.NewBinder()
	echo.HTTPErrorHandler = echoFW.HTTPErrorHandler
	echo.Use(echoFW.LoggerMiddleware(log))
	echo.Use(echoFW.LangMiddleware(config.App.LangDefault))

	return &ServerHTTP{
		echo:   echo,
		config: config,
	}
}

func (s *ServerHTTP) Start() {
	addr := s.config.App.HostPort()
	appName := color.Format(color.GREEN, s.config.App.Name)
	version := color.Format(color.YELLOW, s.config.App.Version)
	fmt.Printf("%s%s version %s\n\n", cmd.AppLogoASCII, appName, version)
	log.Printf("STARTED HTTP SERVER AT %v\n", addr)
	go func() {
		err := s.echo.Start(addr)
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
	err := s.echo.Shutdown(ctxShutdown)
	if err != nil {
		log.Printf("HTTP SERVER ERROR: %s", err.Error())
	}
	log.Println("SUCCESSFULLY SHUTDOWN HTTP SERVER")
	return nil
}
