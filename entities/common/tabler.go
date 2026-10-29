// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package common

import (
	"time"

	collectionPkg "github.com/dedyf5/resik/pkg/collection"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"gorm.io/gorm"
)

type UserPublicIDProvider interface {
	AllUserPublicIDs() []uuidPkg.UUIDV7
}

func UniqueAllUserPublicIDs[
	T any,
	PT interface {
		*T
		UserPublicIDProvider
	},
](items []T) []uuidPkg.UUIDV7 {
	keys := make(map[uuidPkg.UUIDV7]struct{}, len(items)*2)
	var list []uuidPkg.UUIDV7

	for i := range items {
		var pt PT = &items[i]

		for _, id := range pt.AllUserPublicIDs() {
			if _, exists := keys[id]; !exists {
				keys[id] = struct{}{}
				list = append(list, id)
			}
		}
	}

	return list
}

type TableIdentity struct {
	ID       uint64         `json:"-" gorm:"column:id;type:bigint(20);primaryKey;autoIncrement;"`
	PublicID uuidPkg.UUIDV7 `json:"id" gorm:"column:public_id;type:uuid;unique;not null;"`
}

func (t *TableIdentity) BeforeCreate(tx *gorm.DB) (err error) {
	t.PublicID = uuidPkg.NewUUIDV7()
	return
}

type TableIdentifier interface {
	GetTableIdentity() TableIdentity
}

func UniquePublicIDMap[T TableIdentifier](items []T) map[uuidPkg.UUIDV7]T {
	results := make(map[uuidPkg.UUIDV7]T, len(items))
	for _, i := range items {
		if publicID := i.GetTableIdentity().PublicID; publicID != uuidPkg.Nil {
			results[publicID] = i
		}
	}
	return results
}

type TableAuthor struct {
	CreatedAt         time.Time      `json:"created_at" gorm:"column:created_at;type:datetime;not null;"`
	CreatedByPublicID uuidPkg.UUIDV7 `json:"created_by" gorm:"column:created_by_public_id;type:uuid;not null;"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"column:updated_at;type:datetime;not null;"`
	UpdatedByPublicID uuidPkg.UUIDV7 `json:"updated_by" gorm:"column:updated_by_public_id;type:uuid;not null;"`
}

func (t *TableAuthor) AllUserPublicIDs() []uuidPkg.UUIDV7 {
	return []uuidPkg.UUIDV7{t.CreatedByPublicID, t.UpdatedByPublicID}
}

func (t *TableAuthor) UniqueAllUserPublicIDs() []uuidPkg.UUIDV7 {
	return UniqueAllUserPublicIDs([]TableAuthor{*t})
}

type TableAuthors []TableAuthor

func (ts TableAuthors) UniqueCreatedByPublicIDs() []uuidPkg.UUIDV7 {
	return collectionPkg.Unique(ts, func(t TableAuthor) uuidPkg.UUIDV7 {
		return t.CreatedByPublicID
	})
}

func (ts TableAuthors) UniqueUpdatedByPublicIDs() []uuidPkg.UUIDV7 {
	return collectionPkg.Unique(ts, func(t TableAuthor) uuidPkg.UUIDV7 {
		return t.UpdatedByPublicID
	})
}

func (ts TableAuthors) UniqueAllUserPublicIDs() []uuidPkg.UUIDV7 {
	return UniqueAllUserPublicIDs(ts)
}

type TableAuthorer interface {
	GetTableAuthor() TableAuthor
}

func ExtractTableAuthors[T TableAuthorer](items []T) TableAuthors {
	authors := make(TableAuthors, len(items))
	for i, t := range items {
		authors[i] = t.GetTableAuthor()
	}
	return authors
}
