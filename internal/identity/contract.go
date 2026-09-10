// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package identity

import (
	"context"

	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

//go:generate mockgen -source contract.go -package identity -destination ./mock/identity_resolver.go
type IdentityResolver interface {
	Resolve(c context.Context, tableName string, publicID uuidPkg.UUIDV7) (uint64, error)
	ResolveBatch(c context.Context, tableName string, publicIDs []uuidPkg.UUIDV7) ([]uint64, error)
	GetTenantMerchantIDs(c context.Context, userID uint64) ([]uint64, error)
	GetTenantOutletIDs(c context.Context, userID uint64) ([]uint64, error)

	// Core Generic Authorization Engine
	GetResourceIDs(c context.Context, userID uint64, resourceTable string, permissionCode string) ([]uint64, error)
	HasAccessByPublicID(c context.Context, userID uint64, permissionCode string, resourceTable string, publicID uuidPkg.UUIDV7) (bool, error)
	HasAccessByID(c context.Context, userID uint64, permissionCode string, resourceTable string, resourceID uint64) (bool, error)

	// Domain Convenience Helpers
	GetMerchantIDsByPermission(c context.Context, userID uint64, permissionCode string) ([]uint64, error)
	GetOutletIDsByPermission(c context.Context, userID uint64, permissionCode string) ([]uint64, error)
	InvalidateUserAccessMerchant(ctx context.Context, userID uint64) error
	InvalidateUserAccessOutlet(ctx context.Context, userID uint64) error
}
