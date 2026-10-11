// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"time"

	ctx "github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/groupperiod"
	trxParam "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"github.com/dedyf5/resik/utils/datetime"
)

func (o *BranchOmzetGet) ToParam(c *ctx.Ctx, branchPublicID uuidPkg.UUIDV7) (result *trxParam.BranchOmzetGet, err *resPkg.Status) {
	orderStr := "period"
	if o.Order != nil {
		orderStr = o.GetOrder()
	}

	datetimeStart, err := datetime.FromString(o.GetDatetimeStart(), time.RFC3339, c)
	if err != nil {
		return nil, err
	}

	datetimeEnd, err := datetime.FromString(o.GetDatetimeEnd(), time.RFC3339, c)
	if err != nil {
		return nil, err
	}

	return &trxParam.BranchOmzetGet{
		Ctx:            c,
		BranchPublicID: branchPublicID,
		GroupPeriod: groupperiod.GroupPeriod{
			Mode:          groupperiod.Mode(o.GetMode()),
			DatetimeStart: datetimeStart,
			DatetimeEnd:   datetimeEnd,
			Timezone:      o.GetTimezone(),
		},
		Filter: *goku.NewFilter(o.GetSearch(), o.GetPage(), o.GetLimit()),
		Orders: goku.OrdersBuilder(orderStr),
	}, nil
}
