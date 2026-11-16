// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package drivers

import (
	"testing"
)

func TestNewRedisConnection_NilConfig(t *testing.T) {
	client, cleanup, err := NewRedisConnection(nil)
	if err != nil {
		t.Errorf("expected no error for nil config, got %v", err)
	}
	if client != nil {
		t.Errorf("expected nil client for nil config, got %v", err)
	}
	if cleanup == nil {
		t.Error("expected non-nil cleanup function")
	}
	cleanup()
}

func TestNewRedisConnection_WithConfig(t *testing.T) {
	t.Run("Connection failure with invalid host and default pool size", func(t *testing.T) {
		config := &RedisConfig{
			Host:     host,
			Port:     1,
			PoolSize: 0,
		}

		client, cleanup, err := NewRedisConnection(config)
		if err == nil {
			t.Error("expected error connecting to invalid port, got nil")
			if client != nil {
				err = client.Close()
				if err != nil {
					t.Errorf("expected no error closing client on connection failure, got %v", err)
				}
			}
			if cleanup != nil {
				cleanup()
			}
		} else {
			if client != nil {
				t.Errorf("expected nil client on connection failure, got %v", client)
			}
			if cleanup != nil {
				t.Error("expected nil cleanup on connection failure")
			}
		}
	})
}

func TestNewRedisConfig_HostPort(t *testing.T) {
	config := &RedisConfig{
		Host: "127.0.0.1",
		Port: 6379,
	}

	expected := "127.0.0.1:6379"
	if actual := config.HostPort(); actual != expected {
		t.Errorf("expected HostPort %q, got %q", expected, actual)
	}
}
