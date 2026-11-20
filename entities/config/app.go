// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import "github.com/dedyf5/resik/buildinfo"

// App holds core application metadata such as display name, key identifier, and version
type App struct {
	name    string
	nameKey string
	version string
}

// NewApp initializes an App configuration instance, falling back to buildinfo defaults
// when values are empty and overriding the version when tag generation is active.
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

// Name returns the configured application name.
func (a *App) Name() string {
	return a.name
}

// NameKey returns the normalized key identifier for the application.
func (a *App) NameKey() string {
	return a.nameKey
}

// Version returns the resolved application version string.
func (a *App) Version() string {
	return a.version
}
