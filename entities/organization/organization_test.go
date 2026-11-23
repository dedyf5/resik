// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"testing"

	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

func TestOrganization_TableName(t *testing.T) {
	org := Organization{}
	if got := org.TableName(); got != TABLE_NAME {
		t.Errorf("TableName() = %v, want %q", got, TABLE_NAME)
	}
}

func TestOrganization_GetTableAuthor(t *testing.T) {
	createdByPublicID := uuidPkg.NewUUIDV7()
	updatedByPublicID := uuidPkg.NewUUIDV7()

	org := Organization{
		CreatedByPublicID: createdByPublicID,
		UpdatedByPublicID: updatedByPublicID,
		Name:              "Resik HQ",
		Description:       new("Main Organization"),
	}

	got := org.GetTableAuthor()
	if got.CreatedByPublicID == uuidPkg.Nil || got.CreatedByPublicID != createdByPublicID {
		t.Errorf("CreatedByPublicID = %v, want %v", got.CreatedByPublicID, createdByPublicID)
	}
	if got.UpdatedByPublicID == uuidPkg.Nil || got.UpdatedByPublicID != updatedByPublicID {
		t.Errorf("UpdatedByPublicID = %v, want %v", got.UpdatedByPublicID, updatedByPublicID)
	}
}

func TestOrganization_GetTableAuthors(t *testing.T) {
	t.Run("Non-empty slice", func(t *testing.T) {
		author1 := uuidPkg.NewUUIDV7()
		author2 := uuidPkg.NewUUIDV7()

		orgs := Organizations{
			{
				CreatedByPublicID: author1,
				Name:              "Org One",
			},
			{
				CreatedByPublicID: author2,
				Name:              "Org Two",
			},
		}

		authors := orgs.GetTableAuthors()
		if length := len(authors); length != 2 {
			t.Fatalf("length = %d, want 2", length)
		}
		if got := authors[0].CreatedByPublicID; got == uuidPkg.Nil || got != author1 {
			t.Errorf("authors[0].CreatedByPublicID = %v, want %v", got, author1)
		}
		if got := authors[1].CreatedByPublicID; got == uuidPkg.Nil || got != author2 {
			t.Errorf("authors[1].CreatedByPublicID = %v, want %v", got, author2)
		}
	})

	t.Run("Empty slice", func(t *testing.T) {
		var orgs Organizations
		authors := orgs.GetTableAuthors()
		if length := len(authors); length != 0 {
			t.Errorf("length = %d, want 0", length)
		}
	})
}
