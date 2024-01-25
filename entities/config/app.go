// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"fmt"

	"golang.org/x/text/language"
)

type Env string

const (
	EnvDevelopment Env = "development"
	EnvStaging     Env = "staging"
	EnvProduction  Env = "production"
)

type App struct {
	Name        string
	Version     string
	Env         Env
	LangDefault language.Tag
	Host        string
	Port        uint
}

func (a *App) HostPort() string {
	return fmt.Sprintf("%v:%v", a.Host, a.Port)
}
