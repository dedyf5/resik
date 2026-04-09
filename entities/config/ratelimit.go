// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"time"
)

type RateLimitDriver string

const (
	RateLimitDriverMemory RateLimitDriver = "memory"
	RateLimitDriverRedis  RateLimitDriver = "redis"
)

func (d RateLimitDriver) String() string {
	return string(d)
}

type RateLimit struct {
	Driver RateLimitDriver
	Period time.Duration
	Limit  int64
	Prefix string
}

func (rl *RateLimit) FullPrefix(nameKey string) string {
	return nameKey + ":" + rl.Prefix
}
