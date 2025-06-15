// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package health

import "time"

type CheckStatus string

const (
	StatusUp       CheckStatus = "UP"
	StatusDown     CheckStatus = "DOWN"
	StatusDegraded CheckStatus = "DEGRADED"
)

type CheckDetail struct {
	Name   string      `json:"name"`
	Status CheckStatus `json:"status"`
	Error  *string     `json:"error"`
}

type Checker interface {
	Check() CheckDetail
}

type OverallHealthStatus struct {
	OverallStatus CheckStatus   `json:"overall_status"`
	Timestamp     time.Time     `json:"timestamp"`
	Checks        []CheckDetail `json:"checks"`
}
