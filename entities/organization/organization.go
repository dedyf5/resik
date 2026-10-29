// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	commonEntity "github.com/dedyf5/resik/entities/common"
)

const TABLE_NAME = "organizations"

type Organization struct {
	commonEntity.TableIdentity
	commonEntity.TableAuthor
	Name        string  `json:"name" gorm:"column:name;type:varchar(40);not null;"`
	Description *string `json:"description" gorm:"column:description;type:text;null;"`
}

func (o Organization) GetTableAuthor() commonEntity.TableAuthor {
	return o.TableAuthor
}

func (Organization) TableName() string {
	return TABLE_NAME
}

type Organizations []Organization

func (os Organizations) GetTableAuthors() commonEntity.TableAuthors {
	return commonEntity.ExtractTableAuthors(os)
}
