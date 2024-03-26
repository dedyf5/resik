// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64          `json:"id"`
	Name      sql.NullString `json:"name" gorm:"type:varchar(45);default:null"`
	Username  sql.NullString `json:"username" gorm:"column:username;type:varchar(45);default:null"`
	Password  sql.NullString `json:"password" gorm:"type:varchar(225);default:null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	CreatedBy int64          `json:"created_by" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	UpdatedBy int64          `json:"updated_by" gorm:"not null"`
}
