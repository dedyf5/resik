// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import "github.com/dedyf5/resik/build"

type App struct {
	name    string
	nameKey string
	version string
}

func NewApp(name, nameKey, version string) *App {
	if name == "" {
		name = build.AppName
	}
	if nameKey == "" {
		nameKey = build.AppNameKey
	}
	if version == "" {
		version = build.AppVersion
	}

	return &App{
		name:    name,
		nameKey: nameKey,
		version: version,
	}
}

func (a *App) Name() string {
	return a.name
}

func (a *App) NameKey() string {
	return a.nameKey
}

func (a *App) Version() string {
	return a.version
}
