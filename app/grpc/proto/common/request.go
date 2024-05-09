// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package common

import commonEntity "github.com/dedyf5/resik/entities/common"

func (r *Request) ToRequestEntity() commonEntity.Request {
	return commonEntity.Request{
		Lang: r.Lang,
	}
}
