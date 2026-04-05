// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package check

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dedyf5/resik/pkg/response"
	"google.golang.org/grpc/codes"
)

type CheckStatus string

const (
	StatusUp       CheckStatus = "UP"
	StatusDown     CheckStatus = "DOWN"
	StatusDegraded CheckStatus = "DEGRADED"
	StatusDisabled CheckStatus = "DISABLED"
)

type CheckConfig struct {
	Name    string
	Timeout time.Duration
}

type CheckDetail struct {
	Name   string
	Status CheckStatus
	Error  error
}

func (c *CheckDetail) StatusMessage() string {
	return c.Name + "=" + string(c.Status)
}

func (c *CheckDetail) IsHealthy() bool {
	return c.Status == StatusUp || c.Status == StatusDisabled
}

type OverallHealthStatus struct {
	OverallStatus CheckStatus   `json:"overall_status"`
	Timestamp     time.Time     `json:"timestamp"`
	Checks        []CheckDetail `json:"checks"`
}

func NewOverallHealthStatus(timestamp time.Time, checks ...CheckDetail) *OverallHealthStatus {
	overallStatus := StatusUp

	for _, v := range checks {
		if !v.IsHealthy() {
			overallStatus = v.Status
		}

		if overallStatus == StatusDown {
			break
		}
	}

	return &OverallHealthStatus{
		OverallStatus: overallStatus,
		Timestamp:     timestamp,
		Checks:        checks,
	}
}

func (o *OverallHealthStatus) IsHealthy() bool {
	return o.OverallStatus == StatusUp || o.OverallStatus == StatusDisabled
}

func (o *OverallHealthStatus) NotHealthyMessage() *string {
	if o.IsHealthy() {
		return nil
	}

	messages := []string{}
	for _, v := range o.Checks {
		if !v.IsHealthy() {
			messages = append(messages, v.StatusMessage())
		}
	}

	if len(messages) == 0 {
		return nil
	}

	message := strings.Join(messages, ", ")
	return &message
}

func (o *OverallHealthStatus) Error() error {
	if o.IsHealthy() {
		return nil
	}

	errStrings := []string{}
	for _, v := range o.Checks {
		if v.Error != nil {
			errStrings = append(errStrings, fmt.Sprintf("[%s]", v.Error.Error()))
		}
	}

	if len(errStrings) == 0 {
		return nil
	}

	return errors.New(strings.Join(errStrings, ", "))
}

func (o *OverallHealthStatus) HTTPStatusCode() int {
	switch o.OverallStatus {
	case StatusUp, StatusDisabled:
		return http.StatusOK
	case StatusDown, StatusDegraded:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func (o *OverallHealthStatus) GRPCStatusCode() codes.Code {
	return response.HTTPStatusCodeToGRPCCode(o.HTTPStatusCode())
}
