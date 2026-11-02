// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package redis

import (
	"context"

	"github.com/dedyf5/resik/config"
	checkEntity "github.com/dedyf5/resik/entities/check"
	"github.com/redis/go-redis/v9"
)

type CheckRepo struct {
	redis  *redis.Client
	config checkEntity.CheckConfig
}

func New(redisClient *redis.Client, config config.Config) *CheckRepo {
	cfg := checkEntity.CheckConfig{
		Name:    "redis",
		Timeout: 0,
	}

	if redisClient != nil {
		cfg.Timeout = config.Redis.HealthCheckTimeout
	}

	return &CheckRepo{redis: redisClient, config: cfg}
}

func (r *CheckRepo) Check() checkEntity.CheckDetail {
	detail := checkEntity.CheckDetail{
		Name:   r.config.Name,
		Status: checkEntity.StatusDisabled,
	}

	if r.redis == nil {
		return detail
	}

	detail.Status = checkEntity.StatusUp

	c, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()

	if err := r.redis.Ping(c).Err(); err != nil {
		detail.Status = checkEntity.StatusDown
		detail.Error = err
	}

	return detail
}
