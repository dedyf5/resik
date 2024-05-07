// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/dedyf5/resik/cmd"
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/pkg/color"
	"google.golang.org/grpc"
)

type ServerHTTP struct {
	config     config.Config
	listener   net.Listener
	grpcServer *grpc.Server
	router     *Router
}

func newServerHTTP(config config.Config, router *Router) *ServerHTTP {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.App.Port))
	if err != nil {
		log.Fatalf("HTTP SERVER LISTEN ERROR: %s", err.Error())
		return nil
	}
	server := grpc.NewServer()
	return &ServerHTTP{
		config:     config,
		listener:   listener,
		grpcServer: server,
		router:     router,
	}
}

func (s *ServerHTTP) Start() {
	addr := s.config.App.HostPort()
	appName := color.Format(color.GREEN, s.config.App.Name)
	version := color.Format(color.YELLOW, s.config.App.Version)
	log.Printf("%s%s version %s\n\n", cmd.AppLogoASCII, appName, version)
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
