// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package common

import (
	"reflect"
	"testing"
	"time"

	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

type sampleEntity struct {
	Identity TableIdentity
	Author   TableAuthor
}

func (s sampleEntity) GetTableIdentity() TableIdentity {
	return s.Identity
}

func (s sampleEntity) AllUserPublicIDs() []uuidPkg.UUIDV7 {
	return s.Author.AllUserPublicIDs()
}

func (s sampleEntity) GetTableAuthor() TableAuthor {
	return s.Author
}

func TestTableIdentity_BeforeCreate(t *testing.T) {
	id := &TableIdentity{}
	if err := id.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate() returned unexpected error: %v", err)
	}
	if id.PublicID == uuidPkg.Nil {
		t.Error("expected PublicID to be generated, got uuidPkg.Nil")
	}
}

func TestUniquePublicIDMap(t *testing.T) {
	validID1 := uuidPkg.NewUUIDV7()
	validID2 := uuidPkg.NewUUIDV7()

	items := []sampleEntity{
		{Identity: TableIdentity{PublicID: validID1}},
		{Identity: TableIdentity{PublicID: validID1}},
		{Identity: TableIdentity{PublicID: validID2}},
		{Identity: TableIdentity{PublicID: uuidPkg.Nil}},
	}

	got := UniquePublicIDMap(items)
	if length := len(got); length != 2 {
		t.Fatalf("expected map lenght 2, got %d", length)
	}
	if _, ok := got[validID1]; !ok {
		t.Errorf("expected key %v in map", validID1)
	}
	if _, ok := got[validID2]; !ok {
		t.Errorf("expected key %v in map", validID2)
	}
	if _, ok := got[uuidPkg.Nil]; ok {
		t.Error("did not expect uuidPkg.Nil in map")
	}
}

func TestTableAuthorAndTableAuthors(t *testing.T) {
	userA := uuidPkg.NewUUIDV7()
	userB := uuidPkg.NewUUIDV7()
	userC := uuidPkg.NewUUIDV7()
	now := time.Now()

	author1 := TableAuthor{
		CreatedAt:         now,
		CreatedByPublicID: userA,
		UpdatedAt:         now,
		UpdatedByPublicID: userB,
	}
	author2 := TableAuthor{
		CreatedAt:         now,
		CreatedByPublicID: userB,
		UpdatedAt:         now,
		UpdatedByPublicID: userC,
	}

	t.Run("Single TableAuthor IDs", func(t *testing.T) {
		all := author1.AllUserPublicIDs()
		wantAll := []uuidPkg.UUIDV7{userA, userB}
		if !reflect.DeepEqual(all, wantAll) {
			t.Errorf("AllUserPublicIDs() = %v, want %v", all, wantAll)
		}

		sameAuthor := TableAuthor{
			CreatedByPublicID: userA,
			UpdatedByPublicID: userA,
		}
		uniqueSingle := sameAuthor.UniqueAllUserPublicIDs()
		wantUniqueSingle := []uuidPkg.UUIDV7{userA}
		if !reflect.DeepEqual(uniqueSingle, wantUniqueSingle) {
			t.Errorf("UniqueAllUserPublicIDs() = %v, want %v", uniqueSingle, wantUniqueSingle)
		}
	})

	t.Run("TableAuthors slice helpers", func(t *testing.T) {
		authors := TableAuthors{author1, author2}

		gotCreators := authors.UniqueCreatedByPublicIDs()
		wantCreators := []uuidPkg.UUIDV7{userA, userB}
		if !reflect.DeepEqual(gotCreators, wantCreators) {
			t.Errorf("UniqueCreatedByPublicIDs() = %v, want %v", gotCreators, wantCreators)
		}

		gotUpdaters := authors.UniqueUpdatedByPublicIDs()
		wantUpdaters := []uuidPkg.UUIDV7{userB, userC}
		if !reflect.DeepEqual(gotUpdaters, wantUpdaters) {
			t.Errorf("UniqueUpdatedByPublicIDs() = %v, want %v", gotUpdaters, wantUpdaters)
		}

		gotAllUnique := authors.UniqueAllUserPublicIDs()
		wantAllUnique := []uuidPkg.UUIDV7{userA, userB, userC}
		if !reflect.DeepEqual(gotAllUnique, wantAllUnique) {
			t.Errorf("UniqueAllUserPublicIDs() = %v, want %v", gotAllUnique, wantAllUnique)
		}
	})
}

func TestExtractTableAuthorsAndUniqueAllUserPublicIDs(t *testing.T) {
	userA := uuidPkg.NewUUIDV7()
	userB := uuidPkg.NewUUIDV7()

	entities := []sampleEntity{
		{
			Author: TableAuthor{
				CreatedByPublicID: userA,
				UpdatedByPublicID: userB,
			},
		},
		{
			Author: TableAuthor{
				CreatedByPublicID: userB,
				UpdatedByPublicID: userA,
			},
		},
	}

	extracted := ExtractTableAuthors(entities)
	if length := len(extracted); length != 2 {
		t.Fatalf("expected 2 extracted authors, got %d", length)
	}

	uniqueIDs := UniqueAllUserPublicIDs(entities)
	wantUniqueIDs := []uuidPkg.UUIDV7{userA, userB}
	if !reflect.DeepEqual(uniqueIDs, wantUniqueIDs) {
		t.Errorf("UniqueAllUserPublicIDs() = %v, want %v", uniqueIDs, wantUniqueIDs)
	}
}
