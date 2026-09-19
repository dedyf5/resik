// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import "github.com/dedyf5/resik/entities/organization"

func (*OrganizationDelete) ToOrganization(organizationID uint64) *organization.Organization {
	return &organization.Organization{
		ID: organizationID,
	}
}
