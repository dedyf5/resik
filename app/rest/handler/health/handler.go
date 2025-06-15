// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package health

import (
	"net/http"
	"time"

	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	coreHealth "github.com/dedyf5/resik/core/health"
	"github.com/dedyf5/resik/core/health/response"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/datetime"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	log           *logCtx.Log
	fw            echoFW.IEcho
	healthService coreHealth.IService
}

func New(log *logCtx.Log, fw echoFW.IEcho, hs coreHealth.IService) *HealthHandler {
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

	return &resPkg.Status{
		Code:    code,
		Message: statusMsg,
		Data: &response.HealthHealthz{
			AccessedAt: time.Now().Format(datetime.FormatyyyyMMddHHmmss.ToString()),
		},
	}
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
	httpStatus := http.StatusOK
	if status.OverallStatus != coreHealth.StatusUp {
		httpStatus = http.StatusServiceUnavailable
	}

	return &resPkg.Status{
		Code:    httpStatus,
		Message: string(status.OverallStatus),
		Data: &response.HealthReadyz{
			OverallStatus: string(status.OverallStatus),
			AccessedAt:    status.Timestamp.Format(datetime.FormatyyyyMMddHHmmss.ToString()),
			Checks:        response.HealthReadyzCheckFromCheckDetail(status.Checks),
		},
	}
}
