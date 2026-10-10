// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package permission

type PERMISSION_CODE string

const (
	PERMISSION_CODE_ORG_CREATE PERMISSION_CODE = "organization:create"
	PERMISSION_CODE_ORG_READ   PERMISSION_CODE = "organization:read"
	PERMISSION_CODE_ORG_UPDATE PERMISSION_CODE = "organization:update"
	PERMISSION_CODE_ORG_DELETE PERMISSION_CODE = "organization:delete"

	PERMISSION_CODE_TRX_CREATE PERMISSION_CODE = "transaction:create"
	PERMISSION_CODE_TRX_READ   PERMISSION_CODE = "transaction:read"
	PERMISSION_CODE_TRX_UPDATE PERMISSION_CODE = "transaction:update"
	PERMISSION_CODE_TRX_DELETE PERMISSION_CODE = "transaction:delete"
)

func (u PERMISSION_CODE) String() string {
	return string(u)
}
