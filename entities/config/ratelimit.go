// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"time"

	"github.com/dedyf5/resik/build"
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

func (rl *RateLimit) FullPrefix() string {
	return build.AppNameKey + ":" + rl.Prefix
}
