// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package user defines the User entity, user collections, and lookup helpers.
package user

import (
	"time"

	commonEntity "github.com/dedyf5/resik/entities/common"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

// TABLE_NAME specifies the database table name for User entities.
const TABLE_NAME = "users"

// User represents a system user account with profile details, credentials, and audit timestamps.
type User struct {
	commonEntity.TableIdentity
	FullName    string    `json:"full_name" gorm:"column:full_name;type:varchar(45);not null;"`
	Username    string    `json:"username" gorm:"column:username;type:varchar(45);unique;not null;"`
	Password    string    `json:"-" gorm:"column:password;type:varchar(225);not null;"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;type:datetime;not null;"`
	CreatedByID *uint64   `json:"created_by" gorm:"column:created_by_id;null;"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime;not null;"`
	UpdatedByID *uint64   `json:"updated_by" gorm:"column:updated_by_id;null;"`
}

// GetTableIdentity returns the embedded TableIdentity metadata for the User.
func (u User) GetTableIdentity() commonEntity.TableIdentity {
	return u.TableIdentity
}

// TableName returns the database table name for GORM mapping.
func (User) TableName() string {
	return TABLE_NAME
}

// Users represents a slice of user entities.
type Users []User

// UniquePublicIDMap returns a map of User records keyed by their PublicID.
func (u Users) UniquePublicIDMap() (users map[uuidPkg.UUIDV7]User) {
	return commonEntity.UniquePublicIDMap(u)
}
