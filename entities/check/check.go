// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package check defines health-check status constants, check configurations,
// individual check details, and aggregated health status helpers for HTTP and gRPC response.
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

// CheckStatus represents the operational state of a health check or overall system health.
type CheckStatus string

const (
	StatusUp       CheckStatus = "UP"       // StatusUp indicates the component is healthy and operational.
	StatusDown     CheckStatus = "DOWN"     // StatusDown indicates the component is unavailable or failing.
	StatusDegraded CheckStatus = "DEGRADED" // StatusDegraded indicates the component is operating with reduced capacity or issues.
	StatusDisabled CheckStatus = "DISABLED" // StatusDisabled indicates the health check is disabled.
)

// CheckConfig defines configuration parameters for an individual health check.
type CheckConfig struct {
	Name    string
	Timeout time.Duration
}

// CheckDetail holds the execution result, status, and optional error for a single health check.
type CheckDetail struct {
	Name   string
	Status CheckStatus
	Error  error
}

// StatusMessage format the check name and status as "<name>=<status>".
func (c *CheckDetail) StatusMessage() string {
	return c.Name + "=" + string(c.Status)
}

// IsHealthy reports whether the check detail status is considered healthy (UP or DISABLED).
func (c *CheckDetail) IsHealthy() bool {
	return c.Status == StatusUp || c.Status == StatusDisabled
}

// OverallHealthStatus aggregates individual health check results into a timestamped overall health state.
type OverallHealthStatus struct {
	OverallStatus CheckStatus   `json:"overall_status"`
	Timestamp     time.Time     `json:"timestamp"`
	Checks        []CheckDetail `json:"checks"`
}

// NewOverallHealthStatus builds an OverallHealthStatus from timestamp and check details,
// escalating to non-healthy statuses and stopping early once StatusDown is encountered.
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

// IsHealthy reports whether the overall status is considered healthy (UP or DISABLED).
func (o *OverallHealthStatus) IsHealthy() bool {
	return o.OverallStatus == StatusUp || o.OverallStatus == StatusDisabled
}

// NotHealthyMessage returns a comma-separated summary of unhealthy check status messages,
// or nil when the overall status is healthy or has no unhealthy check messages.
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

// Error aggregates bracketed errors from checks into a single error when unhealthy,
// or returns nil if healthy or no check errors are present.
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

// HTTPStatusCode maps OverallHealthStatus to its corresponding HTTP response status code.
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

// GRPCStatusCode maps the overall HTTP status code to its corresponding gRPC status code.
func (o *OverallHealthStatus) GRPCStatusCode() codes.Code {
	return response.HTTPStatusCodeToGRPCCode(o.HTTPStatusCode())
}
