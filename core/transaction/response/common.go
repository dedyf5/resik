// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	commonCore "github.com/dedyf5/resik/core/common"
	userEntity "github.com/dedyf5/resik/entities/user"
)

func UserFromUserEntity(src *userEntity.User) *commonCore.User {
	if src == nil {
		return nil
	}

	return &commonCore.User{
		Id:       src.PublicID.String32(),
		Name:     src.Name,
		Username: src.Username,
	}
}
