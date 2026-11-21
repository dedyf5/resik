// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package permission

// PERMISSION_CODE represents a permission identifier scoped by resource and action.
type PERMISSION_CODE string

const (
	PERMISSION_CODE_ORG_CREATE PERMISSION_CODE = "organization:create" // PERMISSION_CODE_ORG_CREATE allows creating organizations.
	PERMISSION_CODE_ORG_READ   PERMISSION_CODE = "organization:read"   // PERMISSION_CODE_ORG_READ allows reading organizations.
	PERMISSION_CODE_ORG_UPDATE PERMISSION_CODE = "organization:update" // PERMISSION_CODE_ORG_UPDATE allows updating organizations.
	PERMISSION_CODE_ORG_DELETE PERMISSION_CODE = "organization:delete" // PERMISSION_CODE_ORG_DELETE allows deleting organizations.

	PERMISSION_CODE_BRANCH_CREATE PERMISSION_CODE = "branch:create" // PERMISSION_CODE_BRANCH_CREATE allows creating branches.
	PERMISSION_CODE_BRANCH_READ   PERMISSION_CODE = "branch:read"   // PERMISSION_CODE_BRANCH_READ allows reading branches.
	PERMISSION_CODE_BRANCH_UPDATE PERMISSION_CODE = "branch:update" // PERMISSION_CODE_BRANCH_UPDATE allows updating branches.
	PERMISSION_CODE_BRANCH_DELETE PERMISSION_CODE = "branch:delete" // PERMISSION_CODE_BRANCH_DELETE allows deleting branches.

	PERMISSION_CODE_TRX_CREATE PERMISSION_CODE = "transaction:create" // PERMISSION_CODE_TRX_CREATE allows creating transactions.
	PERMISSION_CODE_TRX_READ   PERMISSION_CODE = "transaction:read"   // PERMISSION_CODE_TRX_READ allows reading transactions.
	PERMISSION_CODE_TRX_UPDATE PERMISSION_CODE = "transaction:update" // PERMISSION_CODE_TRX_UPDATE allows updating transactions.
	PERMISSION_CODE_TRX_DELETE PERMISSION_CODE = "transaction:delete" // PERMISSION_CODE_TRX_DELETE allows deleting transactions.
)

// String returns the string representation of the PERMISSION_CODE.
func (u PERMISSION_CODE) String() string {
	return string(u)
}
