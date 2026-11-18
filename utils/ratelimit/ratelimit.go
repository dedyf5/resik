// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package ratelimit provides rate-limiting capabilities and key generation strategies
// for both REST and gRPC application contexs using memory or Redis storage backend.
package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	jwtCtx "github.com/dedyf5/resik/ctx/jwt"
	"github.com/dedyf5/resik/entities/config"
	resPkg "github.com/dedyf5/resik/pkg/response"
	goredis "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	PrefixUserID = "u"  // PrefixUserID defines the Redis/Memory key prefix for user IDs.
	PrefixIP     = "ip" // PrefixIP defines the Redis/Memory key prefix for IP addresses.
)

// Result represents the outcome of a rate-limiting check operation.
type Result struct {
	Limit      int64 // Maximum number of allowed requests within the period.
	Remaining  int64 // Number of remaining allowed requests in the current period.
	Reset      int64 // Unix timestamp when the rate limit quota resets.
	RetryAfter int64 // Duration in seconds to wait before making a new request.
	Reached    bool  // Indicates whether the rate limit has been reached.
}

// Limiter defines the contract for rate-limiting enforcement and client key extraction.
type Limiter interface {
	Take(c context.Context, key string) (*Result, *resPkg.Status) // Take evaluates and records a hit against the rate limit for the given key.
	GetKeyGRPC(c context.Context) (string, *resPkg.Status)        // GetKeyGRPC extracts a unique client rate-limit key from incoming gRPC context.
	GetKeyREST(c context.Context, ip string) string               // GetKeyREST extracts a unique client rate-limit key from REST context IP.
}

// RateLimit implements the Limiter interface utilizing the underlying ulule/limiter engine.
type RateLimit struct {
	instance *limiter.Limiter
}

// NewRateLimiter initialzes and returns a new rate limiter instance configurad with either Redis or Memory storage store.
func NewRateLimiter(appConf config.App, rateLimitConf config.RateLimit, redisClient *goredis.Client) (Limiter, error) {
	rate := limiter.Rate{
		Period: rateLimitConf.Period,
		Limit:  rateLimitConf.Limit,
	}

	var store limiter.Store
	var err error

	if rateLimitConf.Driver == config.RateLimitDriverRedis {
		if redisClient == nil {
			return nil, fmt.Errorf("rate limit driver is set to 'redis' but Redis connection is not initialized; please check your REDIS_* configuration or switch RATE_LIMIT_DRIVER to '%s'", config.RateLimitDriverMemory)
		}

		store, err = redis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
			Prefix: rateLimitConf.FullPrefix(appConf.NameKey()),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create redis store: %w", err)
		}
	} else {
		store = memory.NewStoreWithOptions(limiter.StoreOptions{
			Prefix: rateLimitConf.FullPrefix(appConf.NameKey()),
		})
	}

	return &RateLimit{
		instance: limiter.New(store, rate),
	}, nil
}

// Take evaluates the rate limit consumtion for specific key and returns the updated rate limit Result.
func (rl *RateLimit) Take(c context.Context, key string) (*Result, *resPkg.Status) {
	limiterCtx, err := rl.instance.Get(c, key)
	if err != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, err)
	}

	now := time.Now().Unix()
	retryAfter := max(limiterCtx.Reset-now, 0)

	return &Result{
		Limit:      limiterCtx.Limit,
		Remaining:  limiterCtx.Remaining,
		Reset:      limiterCtx.Reset,
		RetryAfter: retryAfter,
		Reached:    limiterCtx.Reached,
	}, nil
}

// KeyUserID formats and returns a rate limit key string for a given user ID.
func (rl *RateLimit) KeyUserID(userID uint64) string {
	return fmt.Sprintf("%s:%d", PrefixUserID, userID)
}

// KeyIP formats and returns a rate limit key string for a given IP address, wrapping IPv6 in brackets.
func (rl *RateLimit) KeyIP(ip string) string {
	if strings.Contains(ip, ":") {
		return fmt.Sprintf("%s:[%s]", PrefixIP, ip)
	}
	return fmt.Sprintf("%s:%s", PrefixIP, ip)
}

// GetKeyUserID retrieves the formatted user ID rate limit key if authenticated claim exist in the context.
func (rl *RateLimit) GetKeyUserID(c context.Context) string {
	claim := jwtCtx.AuthClaimsFromContext(c)
	if claim != nil && claim.UserID() > 0 {
		return rl.KeyUserID(claim.UserID())
	}
	return ""
}

// GetKeyGRPC resolves the appropriate rate-limit key from a gRPC incoming context, falling back from User ID to Clodflare/X-Forwarded headers or peer address.
func (rl *RateLimit) GetKeyGRPC(c context.Context) (key string, err *resPkg.Status) {
	if userKey := rl.GetKeyUserID(c); userKey != "" {
		return userKey, nil
	}

	if md, ok := metadata.FromIncomingContext(c); ok {
		if cfIP := md.Get("cf-connecting-ip"); len(cfIP) > 0 {
			return rl.KeyIP(cfIP[0]), nil
		}
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			ips := strings.Split(xff[0], ",")
			return rl.KeyIP(strings.TrimSpace(ips[0])), nil
		}
	}

	if p, ok := peer.FromContext(c); ok {
		host, _, err := net.SplitHostPort(p.Addr.String())
		if err != nil {
			return "", resPkg.NewStatusError(http.StatusInternalServerError, err)
		}
		if host == "" {
			return "", resPkg.NewStatusError(http.StatusInternalServerError, errors.New("empty host"))
		}
		return rl.KeyIP(host), nil
	}

	return "", resPkg.NewStatusError(http.StatusInternalServerError, errors.New("client identity not found"))
}

// GetKeyREST resolves the appropriate rate-limit key from a REST context, prioritizing authenticated User ID over client IP.
func (rl *RateLimit) GetKeyREST(c context.Context, ip string) string {
	if userKey := rl.GetKeyUserID(c); userKey != "" {
		return userKey
	}
	return rl.KeyIP(ip)
}
