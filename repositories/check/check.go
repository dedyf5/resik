// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package check

import (
	"database/sql"

	"github.com/dedyf5/resik/config"
	checkEntity "github.com/dedyf5/resik/entities/check"
	"github.com/dedyf5/resik/repositories/check/mysql"
)

//go:generate mockgen -source check.go -package mock -destination ./mock/check.go
type ICheck interface {
	Check() checkEntity.CheckDetail
}

func New(db *sql.DB, config config.Config) ICheck {
	return mysql.New(db, config)
}
