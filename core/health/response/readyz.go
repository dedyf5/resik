// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import "github.com/dedyf5/resik/core/health"

func HealthReadyzCheckFromCheckDetail(src []health.CheckDetail) (res []*HealthReadyzCheck) {
	for _, v := range src {
		res = append(res, &HealthReadyzCheck{
			Name:   v.Name,
			Status: string(v.Status),
			Error:  v.Error,
		})
	}
	return
}
