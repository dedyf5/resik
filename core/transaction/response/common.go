// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	userEntity "github.com/dedyf5/resik/entities/user"
)

func UserFromUserEntity(src *userEntity.User) *User {
	if src == nil {
		return nil
	}

	return &User{
		Id:       src.PublicID.String32(),
		Name:     src.Name,
		Username: src.Username,
	}
}
