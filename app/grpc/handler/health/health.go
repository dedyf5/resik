// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package health

import (
	"context"
	"net/http"
	"time"

	status "github.com/dedyf5/resik/app/grpc/proto/status"
	coreHealth "github.com/dedyf5/resik/core/health"
	"github.com/dedyf5/resik/core/health/response"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/datetime"
	codes "google.golang.org/grpc/codes"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type HealthHandler struct {
	UnimplementedHealthServiceServer
	healthService coreHealth.IService
}

func New(hs coreHealth.IService) *HealthHandler {
	return &HealthHandler{
		healthService: hs,
	}
}

func (h *HealthHandler) HealthzGet(c context.Context, _ *emptypb.Empty) (*HealthHealthzGetRes, error) {
	isLive, statusMsg := h.healthService.LivenessCheck(c)
	if !isLive {
		return nil, &resPkg.Status{
			Code:    http.StatusServiceUnavailable,
			Message: "NOT_SERVING",
		}
	}

	return &HealthHealthzGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: statusMsg,
		},
		Data: &response.HealthHealthz{
			AccessedAt: time.Now().Format(datetime.FormatyyyyMMddHHmmss.ToString()),
		},
	}, nil
}

func (h *HealthHandler) ReadyzGet(c context.Context, _ *emptypb.Empty) (*HealthReadyzGetRes, error) {
	readinessStatus := h.healthService.ReadinessCheck(c)
	protoChecks := make([]*response.HealthReadyzCheck, len(readinessStatus.Checks))
	for i, chk := range readinessStatus.Checks {
		protoChk := &response.HealthReadyzCheck{
			Name:   chk.Name,
			Status: string(chk.Status),
		}
		if chk.Error != nil {
			protoChk.Error = chk.Error
		}
		protoChecks[i] = protoChk
	}

	code := codes.OK
	if readinessStatus.OverallStatus != coreHealth.StatusUp {
		code = codes.Unavailable
	}

	return &HealthReadyzGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(code),
			Message: string(readinessStatus.OverallStatus),
		},
		Data: &response.HealthReadyz{
			OverallStatus: string(readinessStatus.OverallStatus),
			AccessedAt:    readinessStatus.Timestamp.Format(datetime.FormatyyyyMMddHHmmss.ToString()),
			Checks:        protoChecks,
		},
	}, nil
}
