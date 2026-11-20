// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"time"
)

// RateLimitDriver specifies the storage backend used for rate limiting.
type RateLimitDriver string

const (
	RateLimitDriverMemory RateLimitDriver = "memory" // RateLimitDriverMemory stores rate limit counters in local memory.
	RateLimitDriverRedis  RateLimitDriver = "redis"  // RateLimitDriverRedis stores rate limit counters in Redis.
)

// String returns the string representation of the RateLimitDriver.
func (d RateLimitDriver) String() string {
	return string(d)
}

// RateLimit configures rate-limiting driver, time window period, request limit, and key prefix.
type RateLimit struct {
	Driver RateLimitDriver
	Period time.Duration
	Limit  int64
	Prefix string
}

// FullPrefix combines the application name key and rate limit prefix into "<nameKey>:<prefix>".
func (rl *RateLimit) FullPrefix(nameKey string) string {
	return nameKey + ":" + rl.Prefix
}
