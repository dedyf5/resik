// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

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
	PrefixUserID = "u"
	PrefixIP     = "ip"
)

type Result struct {
	Limit      int64
	Remaining  int64
	Reset      int64
	RetryAfter int64
	Reached    bool
}

type Limiter interface {
	Take(c context.Context, key string) (*Result, *resPkg.Status)
	GetKeyGRPC(c context.Context) (string, *resPkg.Status)
	GetKeyREST(c context.Context, ip string) string
}

type RateLimit struct {
	instance *limiter.Limiter
}

func NewRateLimiter(conf config.RateLimit, redisClient *goredis.Client) (Limiter, error) {
	rate := limiter.Rate{
		Period: conf.Period,
		Limit:  conf.Limit,
	}

	var store limiter.Store
	var err error

	if conf.Driver == config.RateLimitDriverRedis {
		if redisClient == nil {
			return nil, fmt.Errorf("rate limit driver is set to 'redis' but Redis connection is not initialized; please check your REDIS_* configuration or switch RATE_LIMIT_DRIVER to '%s'", config.RateLimitDriverMemory)
		}

		store, err = redis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
			Prefix: conf.FullPrefix(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create redis store: %w", err)
		}
	} else {
		store = memory.NewStoreWithOptions(limiter.StoreOptions{
			Prefix: conf.FullPrefix(),
		})
	}

	return &RateLimit{
		instance: limiter.New(store, rate),
	}, nil
}

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

func (rl *RateLimit) KeyUserID(userID uint64) string {
	return fmt.Sprintf("%s:%d", PrefixUserID, userID)
}

func (rl *RateLimit) KeyIP(ip string) string {
	if strings.Contains(ip, ":") {
		return fmt.Sprintf("%s:[%s]", PrefixIP, ip)
	}
	return fmt.Sprintf("%s:%s", PrefixIP, ip)
}

func (rl *RateLimit) GetKeyUserID(c context.Context) string {
	claim := jwtCtx.AuthClaimsFromContext(c)
	if claim != nil && claim.UserID > 0 {
		return rl.KeyUserID(claim.UserID)
	}
	return ""
}

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

func (rl *RateLimit) GetKeyREST(c context.Context, ip string) string {
	if userKey := rl.GetKeyUserID(c); userKey != "" {
		return userKey
	}
	return rl.KeyIP(ip)
}
