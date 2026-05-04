// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"time"

	"github.com/dedyf5/resik/pkg/collection"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"gorm.io/gorm"
)

const TABLE_NAME = "merchants"

type Merchant struct {
	ID          uint64         `json:"-" gorm:"primaryKey;autoIncrement;"`
	PublicID    uuidPkg.UUIDV7 `json:"id" gorm:"column:public_id;type:uuid;unique;not null;"`
	OwnerID     uint64         `json:"owner_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"type:varchar(40);not null"`
	Description *string        `json:"description" gorm:"type:text;null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:datetime;not null;"`
	CreatedBy   uint64         `json:"created_by" gorm:"not null"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:datetime;not null;"`
	UpdatedBy   uint64         `json:"updated_by" gorm:"not null"`
}

func (m *Merchant) BeforeCreate(tx *gorm.DB) (err error) {
	m.PublicID, err = uuidPkg.NewUUIDV7()
	return
}

func (m *Merchant) AllUserIDs() [3]uint64 {
	return [3]uint64{m.OwnerID, m.CreatedBy, m.UpdatedBy}
}

func (m *Merchant) UniqueAllUserIDs() []uint64 {
	keys := make(map[uint64]bool, 3)
	var list []uint64

	ids := m.AllUserIDs()
	for _, id := range ids {
		if !keys[id] {
			keys[id] = true
			list = append(list, id)
		}
	}
	return list
}

func (Merchant) TableName() string {
	return TABLE_NAME
}

type Tabler interface {
	TableName() string
}

type Merchants []Merchant

func (ms Merchants) UniqueOwnerIDs() []uint64 {
	return collection.Unique(ms, func(m Merchant) uint64 {
		return m.OwnerID
	})
}

func (ms Merchants) UniqueCreatedBys() []uint64 {
	return collection.Unique(ms, func(m Merchant) uint64 {
		return m.CreatedBy
	})
}

func (ms Merchants) UniqueUpdatedBys() []uint64 {
	return collection.Unique(ms, func(m Merchant) uint64 {
		return m.UpdatedBy
	})
}

func (ms Merchants) UniqueAllUserIDs() []uint64 {
	keys := make(map[uint64]bool, len(ms)*2)
	var list []uint64

	for _, m := range ms {
		ids := m.AllUserIDs()
		for _, id := range ids {
			if !keys[id] {
				keys[id] = true
				list = append(list, id)
			}
		}
	}
	return list
}
