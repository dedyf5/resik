// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"time"

	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"gorm.io/gorm"
)

const TABLE_NAME = "users"

type User struct {
	ID        uint64         `json:"-" gorm:"primaryKey;autoIncrement;"`
	PublicID  uuidPkg.UUIDV7 `json:"id" gorm:"column:public_id;type:uuid;unique;not null;"`
	Name      string         `json:"name" gorm:"type:varchar(45);not null;"`
	Username  string         `json:"username" gorm:"column:username;type:varchar(45);unique;not null;"`
	Password  string         `json:"-" gorm:"column:password;type:varchar(225);not null;"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:datetime;not null;"`
	CreatedBy *uint64        `json:"created_by" gorm:"null;"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:datetime;not null;"`
	UpdatedBy *uint64        `json:"updated_by" gorm:"null;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.PublicID, err = uuidPkg.NewUUIDV7()
	return
}

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return TABLE_NAME
}

type Users []User

func (u Users) UniqueMap() (users map[uint64]User) {
	users = make(map[uint64]User, len(u))
	for _, v := range u {
		if v.ID != 0 {
			users[v.ID] = v
		}
	}
	return
}
