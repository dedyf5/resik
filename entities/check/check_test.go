// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package check

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
)

func TestCheckDetail(t *testing.T) {
	t.Run("Healthy check details", func(t *testing.T) {
		up := CheckDetail{Name: "db", Status: StatusUp}
		if got, want := up.StatusMessage(), "db=UP"; got != want {
			t.Errorf("StatusMessage() = %q, want %q", got, want)
		}
		if !up.IsHealthy() {
			t.Error("expected StatusUp to be healthy")
		}

		disabled := CheckDetail{Name: "cache", Status: StatusDisabled}
		if !disabled.IsHealthy() {
			t.Error("expected StatusDisabled to be healthy")
		}
	})

	t.Run("Unhealthy check details", func(t *testing.T) {
		degraded := CheckDetail{Name: "queue", Status: StatusDegraded, Error: errors.New("High latency")}
		if degraded.IsHealthy() {
			t.Error("expected StatusDegraded not to be healthy")
		}

		down := CheckDetail{Name: "db", Status: StatusDown, Error: errors.New("connection refused")}
		if down.IsHealthy() {
			t.Error("expected StatusDown not to be healthy")
		}
	})
}

func TestOverallHealthStatus_Healthy(t *testing.T) {
	now := time.Now()
	status := NewOverallHealthStatus(
		now,
		CheckDetail{Name: "db", Status: StatusUp},
		CheckDetail{Name: "redis", Status: StatusDisabled},
	)

	if status.OverallStatus != StatusUp {
		t.Errorf("OverallStatus = %q, want %q", status.OverallStatus, StatusUp)
	}
	if !status.Timestamp.Equal(now) {
		t.Errorf("Timestamp = %v, want %v", status.Timestamp, now)
	}
	if !status.IsHealthy() {
		t.Error("expected OverallHealthStatus to be healthy")
	}
	if msg := status.NotHealthyMessage(); msg != nil {
		t.Errorf("NotHealthyMessage() = %v, want nil", *msg)
	}
	if err := status.Error(); err != nil {
		t.Errorf("Error() = %v, want nil", err)
	}
	if code := status.HTTPStatusCode(); code != http.StatusOK {
		t.Errorf("HTTPStatusCode() = %d, want %d", code, http.StatusOK)
	}
	if grpcCode := status.GRPCStatusCode(); grpcCode != codes.OK {
		t.Errorf("GRPCStatusCode() = %v, want %v", grpcCode, codes.OK)
	}
}

func TestOverallHealthStatus_DegradedAndDown(t *testing.T) {
	now := time.Now()

	t.Run("Degraded status", func(t *testing.T) {
		status := NewOverallHealthStatus(
			now,
			CheckDetail{Name: "db", Status: StatusUp},
			CheckDetail{Name: "search", Status: StatusDegraded, Error: errors.New("slow index")},
		)

		if status.OverallStatus != StatusDegraded {
			t.Errorf("OverallStatus = %q, want %q", status.OverallStatus, StatusDegraded)
		}
		if status.IsHealthy() {
			t.Error("expected OverallHealthStatus not to be healthy")
		}

		msg := status.NotHealthyMessage()
		if expected := "search=DEGRADED"; msg == nil || *msg != expected {
			t.Errorf("NotHealthyMessage() = %v, want %q", msg, expected)
		}

		err := status.Error()
		if expected := "[slow index]"; err == nil || err.Error() != expected {
			t.Errorf("Error() = %v, want %q", err, expected)
		}

		if code := status.HTTPStatusCode(); code != http.StatusServiceUnavailable {
			t.Errorf("HTTPStatusCode() = %d, want %d", code, http.StatusServiceUnavailable)
		}
		if grpcCode := status.GRPCStatusCode(); grpcCode != codes.Unavailable {
			t.Errorf("GRPCStatusCode() = %v, want %v", grpcCode, codes.Unavailable)
		}
	})

	t.Run("Down status stops early and aggregates errors", func(t *testing.T) {
		status := NewOverallHealthStatus(
			now,
			CheckDetail{Name: "db", Status: StatusDown, Error: errors.New("db down")},
			CheckDetail{Name: "cache", Status: StatusDegraded, Error: errors.New("cache slow")},
		)

		if status.OverallStatus != StatusDown {
			t.Errorf("OverallStatus = %q, want %q", status.OverallStatus, StatusDown)
		}

		msg := status.NotHealthyMessage()
		if expected := "db=DOWN, cache=DEGRADED"; msg == nil || *msg != expected {
			t.Errorf("NotHealthyMessage() = %v, want %q", msg, expected)
		}

		err := status.Error()
		if expected := "[db down], [cache slow]"; err == nil || err.Error() != expected {
			t.Errorf("Error() = %v, want %q", err, expected)
		}
	})
}

func TestOverallHealthStatus_EdgeCases(t *testing.T) {
	t.Run("Unhealthy with no checks or errors", func(t *testing.T) {
		status := &OverallHealthStatus{
			OverallStatus: StatusDown,
			Checks:        []CheckDetail{},
		}

		if msg := status.NotHealthyMessage(); msg != nil {
			t.Errorf("NotHealthyMessage() = %v, want nil", *msg)
		}
		if err := status.Error(); err != nil {
			t.Errorf("Error() = %v, want nil", err)
		}
	})

	t.Run("Unknown status returns internal server error", func(t *testing.T) {
		status := &OverallHealthStatus{
			OverallStatus: CheckStatus("UNKNOWN"),
		}

		if code := status.HTTPStatusCode(); code != http.StatusInternalServerError {
			t.Errorf("HTTPStatusCode() = %d, want %d", code, http.StatusInternalServerError)
		}
		if grpcCode := status.GRPCStatusCode(); grpcCode != codes.Internal {
			t.Errorf("GRPCStatusCode() = %v, want %v", grpcCode, codes.Internal)
		}
	})
}
