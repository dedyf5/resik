// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"testing"
	"time"
)

func TestRateLimitDriver_String(t *testing.T) {
	if got := RateLimitDriverMemory.String(); got != "memory" {
		t.Errorf("String() = %q, want %q", got, "memory")
	}
	if got := RateLimitDriverRedis.String(); got != "redis" {
		t.Errorf("String() = %q, want %q", got, "redis")
	}
}

func TestRateLimit_FullPrefix(t *testing.T) {
	rl := &RateLimit{
		Driver: RateLimitDriverMemory,
		Period: time.Minute,
		Limit:  100,
		Prefix: "api_limit",
	}

	got := rl.FullPrefix("resik")
	want := "resik:api_limit"
	if got != want {
		t.Errorf("FullPrefix() = %q, want %q", got, want)
	}
}
