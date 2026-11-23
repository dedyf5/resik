// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package permission

import "testing"

func TestPERMISSION_CODE_String(t *testing.T) {
	tests := []struct {
		name string
		code PERMISSION_CODE
		want string
	}{
		{name: "org create", code: PERMISSION_CODE_ORG_CREATE, want: "organization:create"},
		{name: "branch read", code: PERMISSION_CODE_BRANCH_READ, want: "branch:read"},
		{name: "trx update", code: PERMISSION_CODE_TRX_UPDATE, want: "transaction:update"},
		{name: "trx delete", code: PERMISSION_CODE_TRX_DELETE, want: "transaction:delete"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.code.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}
