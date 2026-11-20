// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import "testing"

func TestModule_HostPort(t *testing.T) {
	m := &Module{
		Host: "localhost",
		Port: 8080,
	}

	got := m.HostPort()
	want := "localhost:8080"
	if got != want {
		t.Errorf("HostPort() = %q, want %q", got, want)
	}
}

func TestModuleType(t *testing.T) {
	t.Run("REST module type", func(t *testing.T) {
		mt := ModuleTypeREST
		if got := mt.String(); got != "REST" {
			t.Errorf("String() = %q, want %q", got, "REST")
		}
		if got := mt.BaseKey(); got != "REST" {
			t.Errorf("BaseKey() = %q, want %q", got, "REST")
		}
		if got := mt.Key("PORT"); got != "REST_PORT" {
			t.Errorf("Key() = %q, want %q", got, "REST_PORT")
		}
	})

	t.Run("gRPC module type", func(t *testing.T) {
		mt := ModuleTypeGRPC
		if got := mt.String(); got != "gRPC" {
			t.Errorf("String() = %q, want %q", got, "gRPC")
		}
		if got := mt.BaseKey(); got != "GRPC" {
			t.Errorf("BaseKey() = %q, want %q", got, "GRPC")
		}
		if got := mt.Key("PORT"); got != "GRPC_PORT" {
			t.Errorf("Key() = %q, want %q", got, "GRPC_PORT")
		}
	})
}

func TestPublic_HostPort(t *testing.T) {
	t.Run("Port is zero return host only", func(t *testing.T) {
		p := &Public{
			Host: "api.example.com",
			Port: 0,
		}
		got := p.HostPort()
		want := p.Host
		if got != want {
			t.Errorf("HostPort() = %q, want %q", got, want)
		}
	})

	t.Run("Port is non-zero returns host and port", func(t *testing.T) {
		p := &Public{
			Host: "api.example.com",
			Port: 9000,
		}
		got := p.HostPort()
		want := "api.example.com:9000"
		if got != want {
			t.Errorf("HostPort() = %q, want %q", got, want)
		}
	})
}
