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
	GetMerchantIDs(c context.Context, userID uint64) ([]uint64, error)
	GetOutletIDs(c context.Context, userID uint64) ([]uint64, error)
	InvalidateUserAccessMerchant(ctx context.Context, userID uint64) error
	InvalidateUserAccessOutlet(ctx context.Context, userID uint64) error
}
