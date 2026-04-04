// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import checkEntity "github.com/dedyf5/resik/entities/check"

func HealthReadyzCheckFromCheckDetail(src []checkEntity.CheckDetail) (res []*HealthReadyzCheck) {
	for _, v := range src {
		var err *string = nil
		if v.Error != nil {
			err = new(v.Error.Error())
		}

		res = append(res, &HealthReadyzCheck{
			Name:   v.Name,
			Status: string(v.Status),
			Error:  err,
		})
	}
	return
}
