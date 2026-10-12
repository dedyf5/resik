// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package identity

import (
	"context"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

//go:generate mockgen -source contract.go -package identity -destination ./mock/identity_resolver.go
type IdentityResolver interface {
	Resolve(c context.Context, tableName string, publicID uuidPkg.UUIDV7) (id uint64, err *resPkg.Status)
	ResolveBatch(c context.Context, tableName string, publicIDs []uuidPkg.UUIDV7) (ids []uint64, err *resPkg.Status)
	GetTenantOrganizationIDs(c context.Context, userID uint64) (ids []uint64, err *resPkg.Status)
	GetTenantBranchIDs(c context.Context, userID uint64) (ids []uint64, err *resPkg.Status)

	// Core Generic Authorization Engine
	GetResourceIDs(c context.Context, userID uint64, resourceTable string, permissionCode string) (ids []uint64, err *resPkg.Status)
	HasAccessByPublicID(c context.Context, lang *langCtx.Lang, userID uint64, permissionCode string, resourceTable string, publicID uuidPkg.UUIDV7) (hasAccess bool, err *resPkg.Status)
	HasAccessByID(c context.Context, lang *langCtx.Lang, userID uint64, permissionCode string, resourceTable string, resourceID uint64) (hasAccess bool, err *resPkg.Status)

	// Domain Convenience Helpers
	GetOrganizationIDsByPermission(c context.Context, userID uint64, permissionCode string) (ids []uint64, err *resPkg.Status)
	GetBranchIDsByPermission(c context.Context, userID uint64, permissionCode string) (ids []uint64, err *resPkg.Status)
	GetBranchPublicIDsByPermission(c context.Context, userID uint64, permissionCode string) (publicIDs []uuidPkg.UUIDV7, err *resPkg.Status)
	InvalidateUserAccessOrganization(ctx context.Context, userID uint64) (err *resPkg.Status)
	InvalidateUserAccessBranch(ctx context.Context, userID uint64) (err *resPkg.Status)
}
