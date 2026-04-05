// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package health

import (
	"net/http"
	"time"

	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	healthCore "github.com/dedyf5/resik/core/health"
	"github.com/dedyf5/resik/core/health/response"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	log           *logCtx.Log
	fw            echoFW.IEcho
	healthService healthCore.IService
}

func New(log *logCtx.Log, fw echoFW.IEcho, hs healthCore.IService) *HealthHandler {
	return &HealthHandler{
		log:           log,
		fw:            fw,
		healthService: hs,
	}
}

// @Summary Liveness check
// @Description Checks if the application is running
// @Tags health
// @Accept json
// @Produce json
// @Param       parameter query commonEntity.Request true "Query Param"
// @Success		200	{object}	resPkg.Response{data=response.HealthHealthz}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     503 {object}	resPkg.Response{data=nil}
// @Router		/healthz [get]
func (h *HealthHandler) HealthHealthzGet(echoCtx echo.Context) error {
	var payload commonEntity.Request

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	isLive, statusMsg := h.healthService.LivenessCheck(echoCtx.Request().Context())
	code := http.StatusOK
	if !isLive {
		statusMsg = "NOT_SERVING"
		code = http.StatusServiceUnavailable
	}

	return resPkg.NewStatusSuccess(
		code,
		statusMsg,
		&response.HealthHealthz{
			AccessedAt: time.Now().UTC().Format(time.RFC3339),
		},
	)
}

// @Summary Readiness check
// @Description Checks if the application and its dependencies are ready
// @Tags health
// @Accept json
// @Produce json
// @Param       parameter query commonEntity.Request true "Query Param"
// @Success		200	{object}	resPkg.Response{data=response.HealthReadyz}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     503 {object}	resPkg.Response{data=nil}
// @Router		/readyz [get]
func (h *HealthHandler) HealthReadyzGet(echoCtx echo.Context) error {
	var payload commonEntity.Request

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	status := h.healthService.ReadinessCheck(echoCtx.Request().Context())

	message := string(status.OverallStatus)
	if msg := status.NotHealthyMessage(); msg != nil {
		message = *msg
	}

	return resPkg.NewStatusMessageData(
		status.HTTPStatusCode(),
		message,
		&response.HealthReadyz{
			OverallStatus: string(status.OverallStatus),
			AccessedAt:    status.Timestamp.UTC().Format(time.RFC3339),
			Checks:        response.HealthReadyzCheckFromCheckDetail(status.Checks),
		},
		status.Error(),
	)
}
