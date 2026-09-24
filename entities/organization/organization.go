// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"time"

	"github.com/dedyf5/resik/pkg/collection"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"gorm.io/gorm"
)

const TABLE_NAME = "organizations"

type Organization struct {
	ID                uint64         `json:"-" gorm:"primaryKey;autoIncrement;"`
	PublicID          uuidPkg.UUIDV7 `json:"id" gorm:"column:public_id;type:uuid;unique;not null;"`
	Name              string         `json:"name" gorm:"type:varchar(40);not null"`
	Description       *string        `json:"description" gorm:"type:text;null"`
	CreatedAt         time.Time      `json:"created_at" gorm:"type:datetime;not null;"`
	CreatedByPublicID uuidPkg.UUIDV7 `json:"created_by" gorm:"column:created_by_public_id;type:uuid;not null;"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"type:datetime;not null;"`
	UpdatedByPublicID uuidPkg.UUIDV7 `json:"updated_by" gorm:"column:updated_by_public_id;type:uuid;not null;"`
}

func (m *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	m.PublicID, err = uuidPkg.NewUUIDV7()
	return
}

func (m *Organization) AllUserPublicIDs() [2]uuidPkg.UUIDV7 {
	return [2]uuidPkg.UUIDV7{m.CreatedByPublicID, m.UpdatedByPublicID}
}

func (m *Organization) UniqueAllUserPublicIDs() []uuidPkg.UUIDV7 {
	keys := make(map[uuidPkg.UUIDV7]bool, 2)
	var list []uuidPkg.UUIDV7

	ids := m.AllUserPublicIDs()
	for _, id := range ids {
		if !keys[id] {
			keys[id] = true
			list = append(list, id)
		}
	}
	return list
}

func (Organization) TableName() string {
	return TABLE_NAME
}

type Tabler interface {
	TableName() string
}

type Organizations []Organization

func (ms Organizations) UniqueCreatedByPublicIDs() []uuidPkg.UUIDV7 {
	return collection.Unique(ms, func(m Organization) uuidPkg.UUIDV7 {
		return m.CreatedByPublicID
	})
}

func (ms Organizations) UniqueUpdatedByPublicIDs() []uuidPkg.UUIDV7 {
	return collection.Unique(ms, func(m Organization) uuidPkg.UUIDV7 {
		return m.UpdatedByPublicID
	})
}

func (ms Organizations) UniqueAllUserPublicIDs() []uuidPkg.UUIDV7 {
	keys := make(map[uuidPkg.UUIDV7]bool, len(ms)*2)
	var list []uuidPkg.UUIDV7

	for _, m := range ms {
		ids := m.AllUserPublicIDs()
		for _, id := range ids {
			if !keys[id] {
				keys[id] = true
				list = append(list, id)
			}
		}
	}
	return list
}
