// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package dto

import (
	"github.com/dedyf5/resik/entities/organization"
	"github.com/dedyf5/resik/entities/user"
)

type Organization struct {
	organization.Organization
	Owner   *user.User `json:"owner"`
	Creator *user.User `json:"creator"`
	Updater *user.User `json:"updater"`
}

func OrganizationFromEntity(data organization.Organization, users map[uint64]user.User) Organization {
	var owner, creator, updater *user.User

	if u, ok := users[data.OwnerID]; ok {
		owner = &u
	}
	if u, ok := users[data.CreatedBy]; ok {
		creator = &u
	}
	if u, ok := users[data.UpdatedBy]; ok {
		updater = &u
	}

	return Organization{
		Organization: data,
		Owner:        owner,
		Creator:      creator,
		Updater:      updater,
	}
}

type Organizations []Organization

func OrganizationsFromEntity(data organization.Organizations, users map[uint64]user.User) Organizations {
	n := len(data)
	if n == 0 {
		return Organizations{}
	}

	res := make(Organizations, n)

	for i, m := range data {
		res[i] = OrganizationFromEntity(m, users)
	}

	return res
}

type OrganizationsResult struct {
	Data  Organizations `json:"data"`
	Total int64         `json:"total"`
}

var (
	OrganizationsResultEmpty = OrganizationsResult{
		Data: Organizations{},
	}
)
