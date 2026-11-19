// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package ratelimit

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	configEntity "github.com/dedyf5/resik/entities/config"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func TestNewRateLimiter(t *testing.T) {
	appConf := configEntity.NewApp("TestApp", "testapp", "0.1")

	t.Run("Success with Memory Driver", func(t *testing.T) {
		rateLimitConf := configEntity.RateLimit{
			Driver: configEntity.RateLimitDriverMemory,
			Period: time.Minute,
			Limit:  10,
		}
		limiter, err := NewRateLimiter(*appConf, rateLimitConf, nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if limiter == nil {
			t.Error("expected non-nil limiter intance")
		}
	})

	t.Run("Error with Redis Driver and nil client", func(t *testing.T) {
		rateLimitConf := configEntity.RateLimit{
			Driver: configEntity.RateLimitDriverRedis,
			Period: time.Minute,
			Limit:  10,
		}
		limiter, err := NewRateLimiter(*appConf, rateLimitConf, nil)
		if err == nil {
			t.Error("expected error when redis client is nil, got nil")
		}
		if limiter != nil {
			t.Error("expected nil limiter on error")
		}
	})
}

func TestRateLimit_Take(t *testing.T) {
	appConf := configEntity.NewApp("TestApp", "testapp", "0.1")
	rateLimitConf := configEntity.RateLimit{
		Driver: configEntity.RateLimitDriverMemory,
		Period: time.Minute,
		Limit:  2,
	}

	limiter, err := NewRateLimiter(*appConf, rateLimitConf, nil)
	if err != nil {
		t.Fatalf("failed to create rate limiter: %v", err)
	}

	c := context.Background()

	t.Run("Take within limit", func(t *testing.T) {
		key := limiter.(*RateLimit).KeyUserID(999)
		res, err := limiter.Take(c, key)
		if err != nil {
			t.Errorf("expected no status error, got %v", err)
		}
		if res.Reached {
			t.Errorf("expected rate limit not reached")
		}
		if res.Remaining != 1 {
			t.Errorf("expected remaning 1, got %d", res.Remaining)
		}
	})

	t.Run("Take exceeding limit", func(t *testing.T) {
		key := limiter.(*RateLimit).KeyUserID(888)
		_, _ = limiter.Take(c, key)
		_, _ = limiter.Take(c, key)
		res, err := limiter.Take(c, key)
		if err != nil {
			t.Errorf("expected no status error, got %v", err)
		}
		if !res.Reached {
			t.Error("expected rate limit to be reached")
		}
	})
}

func TestKeyFormatting(t *testing.T) {
	rl := &RateLimit{}

	t.Run("KeyUserID formatting", func(t *testing.T) {
		if got := rl.KeyUserID(123); got != "u:123" {
			t.Errorf("expected %q, got %q", "u:123", got)
		}
	})

	t.Run("KeyIP formatting IPv4", func(t *testing.T) {
		if got := rl.KeyIP("127.0.0.1"); got != "ip:127.0.0.1" {
			t.Errorf("expected %q, got %q", "ip:127.0.0.1", got)
		}
	})

	t.Run("KeyIP formatting IPv6", func(t *testing.T) {
		if got := rl.KeyIP("2001:db8::1"); got != "ip:[2001:db8::1]" {
			t.Errorf("expected %q, got %q", "ip:[2001:db8::1]", got)
		}
	})
}

func TestGetKeyREST(t *testing.T) {
	rl := &RateLimit{}
	c := context.Background()

	t.Run("GetKeyREST falls back to IP when no user claim", func(t *testing.T) {
		got := rl.GetKeyREST(c, "192.168.1.50")
		if expected := "ip:192.168.1.50"; got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}

func TestGetKeyGRPC(t *testing.T) {
	rl := &RateLimit{}

	t.Run("GetKeyGRPC with cf-connecting-ip metadata", func(t *testing.T) {
		md := metadata.New(map[string]string{
			"cf-connecting-ip": "203.0.113.195",
		})
		c := metadata.NewIncomingContext(context.Background(), md)

		key, err := rl.GetKeyGRPC(c)
		if err != nil {
			t.Errorf("expected no status error, got %v", err)
		}
		if expected := "ip:203.0.113.195"; key != expected {
			t.Errorf("expected %q, got %q", expected, key)
		}
	})

	t.Run("GetKeyGRPC with x-forwarded-for metadata", func(t *testing.T) {
		md := metadata.New(map[string]string{
			"x-forwarded-for": "198.51.100.15, 10.0.0.1",
		})
		c := metadata.NewIncomingContext(context.Background(), md)

		key, err := rl.GetKeyGRPC(c)
		if err != nil {
			t.Errorf("expected no status error, got %v", err)
		}
		if expected := "ip:198.51.100.15"; key != expected {
			t.Errorf("expected %q, got %q", expected, key)
		}
	})

	t.Run("GetKeyGRPC with peer address", func(t *testing.T) {
		tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
		p := &peer.Peer{Addr: tcpAddr}
		c := peer.NewContext(context.Background(), p)

		key, err := rl.GetKeyGRPC(c)
		if err != nil {
			t.Errorf("expected no status error, got %v", err)
		}
		if expected := "ip:127.0.0.1"; key != expected {
			t.Errorf("expected %q, got %q", expected, key)
		}
	})

	t.Run("GetKeyGRPC client indentity not found error", func(t *testing.T) {
		c := context.Background()
		_, err := rl.GetKeyGRPC(c)
		if err == nil {
			t.Error("expected error when identity cannot be resolved, got nil")
		}
		if err.Code != http.StatusInternalServerError {
			t.Errorf("expected code %d, got %d", http.StatusInternalServerError, err.Code)
		}
	})
}
