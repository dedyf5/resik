// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package mysql

import (
	"context"
	"database/sql"

	"github.com/dedyf5/resik/config"
	checkEntity "github.com/dedyf5/resik/entities/check"
)

type CheckRepo struct {
	db     *sql.DB
	config checkEntity.CheckConfig
}

func New(db *sql.DB, config config.Config) *CheckRepo {
	cfg := checkEntity.CheckConfig{
		Name:    "database",
		Timeout: config.Database.HealthCheckTimeout,
	}
	return &CheckRepo{db: db, config: cfg}
}

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
