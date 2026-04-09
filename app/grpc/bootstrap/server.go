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

	"github.com/dedyf5/resik/app/grpc/middleware"
	"github.com/dedyf5/resik/buildinfo"
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/pkg/color"
	"google.golang.org/grpc"
)

type ServerHTTP struct {
	config      config.Config
	listener    net.Listener
	grpcServer  *grpc.Server
	router      *Router
	interceptor *middleware.Interceptor
}

func newServerHTTP(c context.Context, config config.Config, router *Router, interceptor *middleware.Interceptor) *ServerHTTP {
	lc := net.ListenConfig{}
	listener, err := lc.Listen(c, "tcp", fmt.Sprintf(":%d", config.Module.Port))
	if err != nil {
		log.Fatalf("HTTP SERVER LISTEN ERROR: %s", err.Error())
		return nil
	}
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.Unary, interceptor.RateLimit),
	)
	return &ServerHTTP{
		config:      config,
		listener:    listener,
		grpcServer:  server,
		router:      router,
		interceptor: interceptor,
	}
}

func (s *ServerHTTP) Start() {
	appName := color.Format(color.GREEN, s.config.App.Name())
	appVersion := color.Format(color.YELLOW, s.config.App.Version())

	moduleName := color.Format(color.GREEN, s.config.Module.Name)
	moduleVersion := color.Format(color.YELLOW, s.config.Module.Version)

	addr := s.config.Module.HostPort()

	fmt.Printf("%s\n\n", buildinfo.FrameworkLogoASCIIVersion)
	fmt.Printf("%s version %s\n", appName, appVersion)
	fmt.Printf("%s version %s\n\n", moduleName, moduleVersion)
	log.Printf("STARTED HTTP SERVER AT %v\n", addr)

	s.router.routerSetup(s.grpcServer)

	go func() {
		if err := s.grpcServer.Serve(s.listener); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				log.Println("HTTP SERVER CLOSED")
			} else {
				log.Fatalf("HTTP SERVER ERROR: %s", err.Error())
			}
		}
	}()
}

func (s *ServerHTTP) Close() {
	log.Println("SHUTTING DOWN HTTP SERVER")
	s.grpcServer.Stop()
	log.Println("SUCCESSFULLY SHUTDOWN HTTP SERVER")
}
