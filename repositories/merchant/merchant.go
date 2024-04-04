// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"net/http"

	"github.com/dedyf5/resik/ctx"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (r *MerchantRepo) MerchantInsert(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, status *resPkg.Status) {
	result := r.DB.WithContext(ctx.Context).Create(merchant)
	if result.Error != nil {
		return false, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: result.Error,
		}
	}
	return true, nil
}
