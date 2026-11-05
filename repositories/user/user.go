// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/drivers"
	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"github.com/dedyf5/resik/repositories/user/mysql"
	"github.com/dedyf5/resik/repositories/user/postgres"
	"gorm.io/gorm"
)

//go:generate mockgen -source user.go -package mock -destination ./mock/user.go
type IUser interface {
	UserByID(ctx *ctx.Ctx, userID uint64) (user *userEntity.User, err *resPkg.Status)
	UserByUsername(ctx *ctx.Ctx, username string) (user *userEntity.User, err *resPkg.Status)
	UsersGetByPublicIDs(ctx *ctx.Ctx, userIDs []uuidPkg.UUIDV7) (users userEntity.Users, err *resPkg.Status)
}

func New(db *gorm.DB, engine drivers.SQLEngine) IUser {
	switch engine {
	case drivers.MySQL:
		return mysql.New(db)
	case drivers.PostgreSQL:
		return postgres.New(db)
	}
	return nil
}
