// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package branch

import (
	"net/http"
	"testing"

	commonEntity "github.com/dedyf5/resik/entities/common"
	"github.com/dedyf5/resik/entities/organization"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

func TestBranch_TableName(t *testing.T) {
	b := Branch{}
	if got := b.TableName(); got != TABLE_NAME {
		t.Errorf("TableName() = %q, want %q", got, TABLE_NAME)
	}
}

func TestBranch_IsBelongToOrganization(t *testing.T) {
	orgID := uuidPkg.NewUUIDV7()
	otherOrgID := uuidPkg.NewUUIDV7()

	b := &Branch{
		Name: "Main Branch",
		Organization: organization.Organization{
			TableIdentity: commonEntity.TableIdentity{
				PublicID: orgID,
			},
		},
	}

	t.Run("Matching organization public ID", func(t *testing.T) {
		ok, err := b.IsBelongToOrganization(&orgID)
		if !ok {
			t.Error("expected IsBelongToOrganization() to return true, got false")
		}
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("Mismatched organization public ID returns unprocessable entity", func(t *testing.T) {
		ok, err := b.IsBelongToOrganization(&otherOrgID)
		if ok {
			t.Error("expected IsBelongToOrganization() to return false, got true")
		}
		if err == nil {
			t.Fatal("expected status error, got nil")
		}
		if err.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected error code %d, got %d", err.Code, http.StatusUnprocessableEntity)
		}
	})
}
