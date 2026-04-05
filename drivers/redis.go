// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package drivers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host               string
	Port               int
	Username           string
	Password           string
	Database           int
	PoolSize           int
	HealthCheckTimeout time.Duration
}

func NewRedisConnection(config *RedisConfig) (*redis.Client, func(), error) {
	if config == nil {
		return nil, func() {}, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.HostPort(),
		Username: config.Username,
		Password: config.Password,
		DB:       config.Database,
		PoolSize: config.PoolSize,
	})

	if config.PoolSize == 0 {
		client.Options().PoolSize = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	cleanup := func() {
		if err := client.Close(); err != nil {
			log.Printf("failed to close redis connection: %v", err)
		}
	}

	return client, cleanup, nil
}

func (r *RedisConfig) HostPort() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
