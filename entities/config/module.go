// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

// Env represents the runtime environment tier of a module.
type Env string

const (
	EnvDevelopment Env = "development" // EnvDevelopment identifies a local or development environment.
	EnvStaging     Env = "staging"     // EnvStaging identifies a pre-production staging environment.
	EnvProduction  Env = "production"  // EnvProduction identifies a live production environment.
)

// ModuleType identifies the transport protocol exposed by a module.
type ModuleType string

const (
	ModuleTypeREST ModuleType = "REST" // ModuleTypeREST represents an HTTP REST service module.
	ModuleTypeGRPC ModuleType = "gRPC" // ModuleTypeGRPC represents a gRPC service module.
)

// Module describes a service module's identity, transport type, environment, locale, network binding, and public endpoint.
type Module struct {
	Name        string
	NameKey     string
	Type        ModuleType
	Env         Env
	LangDefault language.Tag
	Host        string
	Port        uint
	Public      Public
}

// Public defines externally accessable host, port, schema, and base path settings for a module.
type Public struct {
	Host     string
	Port     uint
	Schema   string
	BasePath string
}

// HostPort formats the module's internal bind address as "host:port".
func (m *Module) HostPort() string {
	return fmt.Sprintf("%v:%v", m.Host, m.Port)
}

// String returns the string representation of the ModuleType.
func (t ModuleType) String() string {
	return string(t)
}

// BaseKey return the uppercase prefix key derived from the ModuleType.
func (t ModuleType) BaseKey() string {
	return strings.ToUpper(string(t))
}

// Key formats a module-scoped configuration key using the module's uppercase base key.
func (t ModuleType) Key(k string) string {
	return fmt.Sprintf("%v_%v", t.BaseKey(), k)
}

// HostPort returns the public host address, appending ":port" only when Port is non-zero.
func (p *Public) HostPort() string {
	if p.Port == 0 {
		return p.Host
	}
	return fmt.Sprintf("%v:%v", p.Host, p.Port)
}
