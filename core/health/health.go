// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package health

import (
	"context"

	checkEntity "github.com/dedyf5/resik/entities/check"
)

//go:generate mockgen -source health.go -package mock -destination ./mock/health.go
type IService interface {
	LivenessCheck(c context.Context) (isLive bool, statusMessage string)
	ReadinessCheck(c context.Context) *checkEntity.OverallHealthStatus
}
