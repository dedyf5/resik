// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

type Env string

const (
	EnvDevelopment Env = "development"
	EnvStaging     Env = "staging"
	EnvProduction  Env = "production"
)

type Module string

const (
	ModuleREST Module = "REST"
	ModuleGRPC Module = "GRPC"
)

type App struct {
	Name        string
	Version     string
	Module      Module
	Env         Env
	LangDefault language.Tag
	Host        string
	Port        uint
}

func (a *App) HostPort() string {
	return fmt.Sprintf("%v:%v", a.Host, a.Port)
}

func (a *App) APIDocDescription() string {
	return fmt.Sprintf("%v API Documentation", a.Name)
}

func (t Module) DirectoryName() string {
	return strings.ToLower(string(t))
}
