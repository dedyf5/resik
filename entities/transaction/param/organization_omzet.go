// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package param

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/groupperiod"
	"github.com/dedyf5/resik/pkg/goku"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

type OrganizationOmzetGet struct {
	Ctx                  *ctx.Ctx
	OrganizationPublicID uuidPkg.UUIDV7
	GroupPeriod          groupperiod.GroupPeriod
	Filter               goku.Filter
	Orders               []goku.Order
}
