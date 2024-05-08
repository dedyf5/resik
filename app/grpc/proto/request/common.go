// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import commonEntity "github.com/dedyf5/resik/entities/common"

func (c *Common) ToRequestEntity() commonEntity.Request {
	return commonEntity.Request{
		Lang: c.Lang,
	}
}
