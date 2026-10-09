// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"errors"
	"net/http"

	"github.com/dedyf5/resik/ctx"
	branchEntity "github.com/dedyf5/resik/entities/branch"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"gorm.io/gorm"
)

func (r *OrganizationRepo) BranchGetByPublicID(ctx *ctx.Ctx, branchPublicID *uuidPkg.UUIDV7) (branch *branchEntity.Branch, err *resPkg.Status) {
	query := r.DB.WithContext(ctx.Context).
		Table(branchEntity.TABLE_NAME).
		Preload("Organization").
		Where("public_id = ?", branchPublicID).
		First(&branch)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, query.Error)
	}
	return
}
