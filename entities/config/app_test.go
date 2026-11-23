// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"testing"

	"github.com/dedyf5/resik/buildinfo"
)

func TestNewApp(t *testing.T) {
	t.Run("Custom values", func(t *testing.T) {
		app := NewApp("CustomApp", "custom_app", "1.2.3")
		if want := "CustomApp"; app.Name() != want {
			t.Errorf("Name() = %q, want %q", app.Name(), want)
		}
		if want := "custom_app"; app.NameKey() != want {
			t.Errorf("NameKey() = %q, want %q", app.NameKey(), want)
		}

		wantVersion := "1.2.3"
		if buildinfo.GetAppVersionGenerator() == buildinfo.VersionGeneratorTag {
			wantVersion = buildinfo.AppVersion
		}
		if app.Version() != wantVersion {
			t.Errorf("Version() = %q, want %q", app.Version(), wantVersion)
		}
	})

	t.Run("Empty values fallback to buildinfo defaults", func(t *testing.T) {
		app := NewApp("", "", "")
		if app.Name() != buildinfo.FrameworkName {
			t.Errorf("Name() = %q, want %q", app.Name(), buildinfo.FrameworkName)
		}
		if app.NameKey() != buildinfo.FrameworkNameKey {
			t.Errorf("NameKey() = %q, want %q", app.NameKey(), buildinfo.FrameworkNameKey)
		}

		wantVersion := buildinfo.FrameworkVersion
		if buildinfo.GetAppVersionGenerator() == buildinfo.VersionGeneratorTag {
			wantVersion = buildinfo.AppVersion
		}
		if app.Version() != wantVersion {
			t.Errorf("Version() = %q, want %q", app.Version(), wantVersion)
		}
	})
}
