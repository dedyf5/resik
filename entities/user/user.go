// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"time"
)

type User struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name" gorm:"type:varchar(45)"`
	Username  string    `json:"username" gorm:"column:username;type:varchar(45);not null"`
	Password  string    `json:"password" gorm:"type:varchar(225)"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	CreatedBy uint64    `json:"created_by" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	UpdatedBy uint64    `json:"updated_by" gorm:"not null"`
}

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return "user"
}
