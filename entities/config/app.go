// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import "github.com/dedyf5/resik/buildinfo"

type App struct {
	name    string
	nameKey string
	version string
}

func NewApp(name, nameKey, version string) *App {
	if name == "" {
		name = buildinfo.FrameworkName
	}
	if nameKey == "" {
		nameKey = buildinfo.FrameworkNameKey
	}
	if version == "" {
		version = buildinfo.FrameworkVersion
	}

	if buildinfo.GetAppVersionGenerator() == buildinfo.VersionGeneratorTag {
		version = buildinfo.AppVersion
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
