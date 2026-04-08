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

type ModuleType string

const (
	ModuleTypeREST ModuleType = "REST"
	ModuleTypeGRPC ModuleType = "gRPC"
)

type Module struct {
	Name        string
	NameKey     string
	Version     string
	Type        ModuleType
	Env         Env
	LangDefault language.Tag
	Host        string
	Port        uint
	Public      Public
}

type Public struct {
	Host     string
	Port     uint
	Schema   string
	BasePath string
}

func (a *Module) HostPort() string {
	return fmt.Sprintf("%v:%v", a.Host, a.Port)
}

func (a *Module) APIDocDescription() string {
	return fmt.Sprintf("%v API Documentation", a.Name)
}

func (t ModuleType) String() string {
	return string(t)
}

func (t ModuleType) BaseKey() string {
	return strings.ToUpper(string(t))
}

func (t ModuleType) Key(k string) string {
	return fmt.Sprintf("%v_%v", t.BaseKey(), k)
}

func (p *Public) HostPort() string {
	if p.Port == 0 {
		return p.Host
	}
	return fmt.Sprintf("%v:%v", p.Host, p.Port)
}
