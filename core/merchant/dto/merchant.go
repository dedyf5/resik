// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package dto

import (
	"github.com/dedyf5/resik/entities/merchant"
	"github.com/dedyf5/resik/entities/user"
)

type Merchant struct {
	merchant.Merchant
	Owner   *user.User `json:"owner"`
	Creator *user.User `json:"creator"`
	Updater *user.User `json:"updater"`
}

func MerchantFromEntity(data merchant.Merchant, users map[uint64]user.User) Merchant {
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

	return Merchant{
		Merchant: data,
		Owner:    owner,
		Creator:  creator,
		Updater:  updater,
	}
}

type Merchants []Merchant

func MerchantsFromEntity(data merchant.Merchants, users map[uint64]user.User) Merchants {
	n := len(data)
	if n == 0 {
		return Merchants{}
	}

	res := make(Merchants, n)

	for i, m := range data {
		res[i] = MerchantFromEntity(m, users)
	}

	return res
}

type MerchantsResult struct {
	Data  Merchants `json:"data"`
	Total int64     `json:"total"`
}

var (
	MerchatsResultEmpty = MerchantsResult{
		Data: Merchants{},
	}
)
