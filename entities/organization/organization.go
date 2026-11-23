// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package organization defines the Organization entity, collection type, table mapping,
// and audit author extraction helpers.
package organization

import (
	commonEntity "github.com/dedyf5/resik/entities/common"
)

// TABLE_NAME specifies the database table name for Organization entities.
const TABLE_NAME = "organizations"

// Organization represents an organization entity with identity, author metadata, and optional description.
type Organization struct {
	commonEntity.TableIdentity
	commonEntity.TableAuthor
	Name        string  `json:"name" gorm:"column:name;type:varchar(40);not null;"`
	Description *string `json:"description" gorm:"column:description;type:text;null;"`
}

// GetTableAuthor returns the embedded TableAuthor audit metadata for the Organization.
func (o Organization) GetTableAuthor() commonEntity.TableAuthor {
	return o.TableAuthor
}

// TableName returns the database table name for GORM mapping.
func (Organization) TableName() string {
	return TABLE_NAME
}

// Organizations represents a slice of Organization entities.
type Organizations []Organization

// GetTableAuthors extracts TableAuthor metadata from all organizations in the slice.
func (os Organizations) GetTableAuthors() commonEntity.TableAuthors {
	return commonEntity.ExtractTableAuthors(os)
}
