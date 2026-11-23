// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"testing"
	"time"

	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

func TestUser_TableName(t *testing.T) {
	u := User{}
	if got := u.TableName(); got != TABLE_NAME {
		t.Errorf("TableName() = %q, want %q", got, TABLE_NAME)
	}
}

func TestUser_GetTableIdentity(t *testing.T) {
	publicID := uuidPkg.NewUUIDV7()
	now := time.Now().UTC()
	createdByID := uint64(1)
	updatedByID := uint64(2)

	u := User{
		PublicID:    publicID,
		FullName:    "Budi Santoso",
		Username:    "budisantoso",
		Password:    "hashed-secret",
		CreatedAt:   now,
		CreatedByID: &createdByID,
		UpdatedAt:   now,
		UpdatedByID: &updatedByID,
	}

	got := u.GetTableIdentity()
	if got.PublicID != publicID {
		t.Errorf("GetTableIdentity().PublicID = %v, want %v", got.PublicID, publicID)
	}
}

func TestUsers_UniquePublicIDMap(t *testing.T) {
	t.Run("Non-empty slice", func(t *testing.T) {
		id1 := uuidPkg.NewUUIDV7()
		id2 := uuidPkg.NewUUIDV7()

		users := Users{
			{
				PublicID: id1,
				FullName: "User One",
				Username: "userone",
			},
			{
				PublicID: id2,
				FullName: "User Two",
				Username: "usertwo",
			},
		}

		mapped := users.UniquePublicIDMap()
		if length := len(mapped); length != 2 {
			t.Fatalf("length = %d, want 2", length)
		}
		if got, expected := mapped[id1].Username, "userone"; got != expected {
			t.Errorf("mapped[%v].Username = %q, want %q", id1, got, expected)
		}
		if got, expected := mapped[id2].Username, "usertwo"; got != expected {
			t.Errorf("mapped[%v].Username = %q, want %q", id2, got, expected)
		}
	})

	t.Run("Empty slice", func(t *testing.T) {
		var users Users
		mapped := users.UniquePublicIDMap()
		if length := len(mapped); length != 0 {
			t.Errorf("length = %d, want 0", length)
		}
	})
}
