// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package health

import "context"

type IService interface {
	LivenessCheck(c context.Context) (isLive bool, statusMessage string)
	ReadinessCheck(c context.Context) OverallHealthStatus
}
