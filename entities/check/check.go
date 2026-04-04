// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package check

import "time"

type CheckStatus string

const (
	StatusUp       CheckStatus = "UP"
	StatusDown     CheckStatus = "DOWN"
	StatusDegraded CheckStatus = "DEGRADED"
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

type OverallHealthStatus struct {
	OverallStatus CheckStatus   `json:"overall_status"`
	Timestamp     time.Time     `json:"timestamp"`
	Checks        []CheckDetail `json:"checks"`
}
