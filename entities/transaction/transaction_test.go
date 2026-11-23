// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"testing"
	"time"

	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

func TestTransaction_TableName(t *testing.T) {
	trx := Transaction{}
	if got := trx.TableName(); got != TABLE_NAME {
		t.Errorf("TableName() = %q, want %q", got, TABLE_NAME)
	}
}

func TestTransaction_GetTableAuthor(t *testing.T) {
	createdByPublicID := uuidPkg.NewUUIDV7()
	updatedByPublicID := uuidPkg.NewUUIDV7()

	trx := Transaction{
		CreatedByPublicID:    createdByPublicID,
		UpdatedByPublicID:    updatedByPublicID,
		OrganizationPublicID: uuidPkg.NewUUIDV7(),
		OrganizationName:     "Resik HQ",
		BranchPublicID:       uuidPkg.NewUUIDV7(),
		BranchName:           "Main Branch",
		TransactedAt:         time.Now().UTC(),
	}

	tableAuthor := trx.GetTableAuthor()
	if got := tableAuthor.CreatedByPublicID; got == uuidPkg.Nil || got != createdByPublicID {
		t.Errorf("CreatedByPublicID = %v, want %v", got, createdByPublicID)
	}
	if got := tableAuthor.UpdatedByPublicID; got == uuidPkg.Nil || got != updatedByPublicID {
		t.Errorf("UpdatedByPublicID = %v, want %v", got, updatedByPublicID)
	}
}

func TestTransactionsGetTableAuthors(t *testing.T) {
	t.Run("Non-empty slice", func(t *testing.T) {
		author1 := uuidPkg.NewUUIDV7()
		author2 := uuidPkg.NewUUIDV7()

		trxs := Transactions{
			{
				CreatedByPublicID: author1,
				OrganizationName:  "Org One",
				BranchName:        "Branch One",
				BillTotal:         50000,
			},
			{
				CreatedByPublicID: author2,
				OrganizationName:  "Org Two",
				BranchName:        "Branch Two",
				BillTotal:         75000,
			},
		}

		authors := trxs.GetTableAuthors()
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
}
