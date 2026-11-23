// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package param defines query and filter parameter structs for organization operations.
package param

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/pkg/goku"
)

// OrganizationsGet holds execution context, organization ID filters, query filters,
// and ordering rules for retrieving organizations.
type OrganizationsGet struct {
	Ctx             *ctx.Ctx
	OrganizationIDs []uint64
	Filter          goku.Filter
	Orders          []goku.Order
}
