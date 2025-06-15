// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package checkers

import (
	"context"
	"database/sql"
	"time"

	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/core/health"
)

type DatabaseChecker struct {
	db  *sql.DB
	cfg DBCheckerConfig
}

type DBCheckerConfig struct {
	Name    string
	Timeout time.Duration
}

func NewDatabaseChecker(db *sql.DB, config config.Config) health.Checker {
	cfg := DBCheckerConfig{Name: "database"}
	if config.Database.HealthCheckTimeoutSeconds > 0 {
		cfg.Timeout = time.Duration(config.Database.HealthCheckTimeoutSeconds) * time.Second
	} else {
		cfg.Timeout = 2 * time.Second
	}
	return &DatabaseChecker{db: db, cfg: cfg}
}

func (dc *DatabaseChecker) Check() health.CheckDetail {
	detail := health.CheckDetail{
		Name:   dc.cfg.Name,
		Status: health.StatusUp,
	}

	ctx, cancel := context.WithTimeout(context.Background(), dc.cfg.Timeout)
	defer cancel()

	if err := dc.db.PingContext(ctx); err != nil {
		errMesssage := err.Error()
		detail.Status = health.StatusDown
		detail.Error = &errMesssage
	}
	return detail
}
