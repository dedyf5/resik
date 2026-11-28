// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package mysql provides a MySQL implementation of the health check repository
// using database connectivity pings.
package mysql

import (
	"context"
	"database/sql"

	"github.com/dedyf5/resik/config"
	checkEntity "github.com/dedyf5/resik/entities/check"
)

// CheckRepo checks MySQL database health using a configured ping timeout.
type CheckRepo struct {
	db     *sql.DB
	config checkEntity.CheckConfig
}

// New creates a new CheckRepo configured with the database health check name
// and timeout from the application configuration.
func New(db *sql.DB, config config.Config) *CheckRepo {
	cfg := checkEntity.CheckConfig{
		Name:    "database",
		Timeout: config.Database.HealthCheckTimeout,
	}
	return &CheckRepo{db: db, config: cfg}
}

// Check pings the database within the configured timeout and returns the
// resulting health check status and any connection error.
func (dc *CheckRepo) Check() checkEntity.CheckDetail {
	detail := checkEntity.CheckDetail{
		Name:   dc.config.Name,
		Status: checkEntity.StatusUp,
	}

	c, cancel := context.WithTimeout(context.Background(), dc.config.Timeout)
	defer cancel()

	if err := dc.db.PingContext(c); err != nil {
		detail.Status = checkEntity.StatusDown
		detail.Error = err
	}

	return detail
}
