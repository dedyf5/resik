// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package health

import (
	"context"
	"net/http"

	status "github.com/dedyf5/resik/app/grpc/proto/status"
	healthCore "github.com/dedyf5/resik/core/health"
	"github.com/dedyf5/resik/core/health/response"
	resPkg "github.com/dedyf5/resik/pkg/response"
	codes "google.golang.org/grpc/codes"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HealthHandler struct {
	UnimplementedHealthServiceServer
	healthService healthCore.IService
}

func New(hs healthCore.IService) *HealthHandler {
	return &HealthHandler{
		healthService: hs,
	}
}

func (h *HealthHandler) HealthzGet(c context.Context, _ *emptypb.Empty) (*HealthHealthzGetRes, error) {
	isLive, statusMsg := h.healthService.LivenessCheck(c)
	if !isLive {
		return nil, resPkg.NewStatusMessage(http.StatusServiceUnavailable, "NOT_SERVING", nil)
	}

	return &HealthHealthzGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: statusMsg,
		},
		Data: &response.HealthHealthz{
			AccessedAt: timestamppb.Now(),
		},
	}, nil
}

func (h *HealthHandler) ReadyzGet(c context.Context, _ *emptypb.Empty) (*HealthReadyzGetRes, error) {
	readinessStatus := h.healthService.ReadinessCheck(c)

	if msg := readinessStatus.NotHealthyMessage(); msg != nil {
		return nil, resPkg.NewStatusMessage(
			readinessStatus.HTTPStatusCode(),
			*msg,
			readinessStatus.Error(),
		)
	}

	return &HealthReadyzGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(readinessStatus.GRPCStatusCode()),
			Message: string(readinessStatus.OverallStatus),
		},
		Data: &response.HealthReadyz{
			OverallStatus: string(readinessStatus.OverallStatus),
			AccessedAt:    timestamppb.New(readinessStatus.Timestamp),
			Checks:        response.HealthReadyzCheckFromCheckDetail(readinessStatus.Checks),
		},
	}, nil
}
