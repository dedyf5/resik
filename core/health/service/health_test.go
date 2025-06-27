// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"context"
	"testing"
	"time"

	"github.com/dedyf5/resik/core/health"
	healthMock "github.com/dedyf5/resik/core/health/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// setup initializes gomock controller, mock checkers, and the health service.
// It returns the slice of mock checkers and the health service instance.
func setup(ctrl *gomock.Controller, numCheckers int) ([]*healthMock.MockChecker, health.IService) {
	mockCheckers := make([]*healthMock.MockChecker, numCheckers)
	healthCheckers := make([]health.Checker, numCheckers)
	for i := range numCheckers {
		mockCheckers[i] = healthMock.NewMockChecker(ctrl)
		healthCheckers[i] = mockCheckers[i]
	}
	healthService := New(healthCheckers)
	return mockCheckers, healthService
}

func TestLivenessCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// No checkers are needed for LivenessCheck as it's a static response.
	// We pass 0 for numCheckers.
	_, healthService := setup(ctrl, 0)

	t.Run("LivenessCheck always returns true and SERVING", func(t *testing.T) {
		isLive, statusMessage := healthService.LivenessCheck(context.Background())
		assert.True(t, isLive)
		assert.Equal(t, "SERVING", statusMessage)
	})
}

func TestReadinessCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	t.Run("All checkers are UP", func(t *testing.T) {
		mockCheckers, healthService := setup(ctrl, 3)
		expectedDetails := []health.CheckDetail{
			{Name: "Checker1", Status: health.StatusUp, Error: nil},
			{Name: "Checker2", Status: health.StatusUp, Error: nil},
			{Name: "Checker3", Status: health.StatusUp, Error: nil},
		}

		// Set expectations for each mock checker
		mockCheckers[0].EXPECT().Check().Return(expectedDetails[0]).Times(1)
		mockCheckers[1].EXPECT().Check().Return(expectedDetails[1]).Times(1)
		mockCheckers[2].EXPECT().Check().Return(expectedDetails[2]).Times(1)

		overallStatus := healthService.ReadinessCheck(ctx)

		assert.Equal(t, health.StatusUp, overallStatus.OverallStatus)
		assert.Len(t, overallStatus.Checks, 3)
		// Assert individual check details. Order is not guaranteed due to goroutines.
		// We use assert.Contains to check if each expected detail is present.
		for _, expected := range expectedDetails {
			assert.Contains(t, overallStatus.Checks, expected)
		}
		// Check if the timestamp is recent (within a few seconds)
		assert.WithinDuration(t, time.Now(), overallStatus.Timestamp, 5*time.Second)
	})

	t.Run("One checker is DOWN", func(t *testing.T) {
		mockCheckers, healthService := setup(ctrl, 2)
		errStr := "database connection failed" // Error message for the DOWN status
		expectedDetails := []health.CheckDetail{
			{Name: "DBChecker", Status: health.StatusDown, Error: &errStr}, // This one is DOWN
			{Name: "APIChecker", Status: health.StatusUp, Error: nil},
		}

		mockCheckers[0].EXPECT().Check().Return(expectedDetails[0]).Times(1)
		mockCheckers[1].EXPECT().Check().Return(expectedDetails[1]).Times(1)

		overallStatus := healthService.ReadinessCheck(ctx)

		assert.Equal(t, health.StatusDown, overallStatus.OverallStatus) // Overall status should be DOWN
		assert.Len(t, overallStatus.Checks, 2)
		for _, expected := range expectedDetails {
			assert.Contains(t, overallStatus.Checks, expected)
		}
	})

	t.Run("One checker is DEGRADED, others UP", func(t *testing.T) {
		mockCheckers, healthService := setup(ctrl, 2)
		errStr := "cache performance degraded"
		expectedDetails := []health.CheckDetail{
			{Name: "CacheChecker", Status: health.StatusDegraded, Error: &errStr}, // This one is DEGRADED
			{Name: "QueueChecker", Status: health.StatusUp, Error: nil},
		}

		mockCheckers[0].EXPECT().Check().Return(expectedDetails[0]).Times(1)
		mockCheckers[1].EXPECT().Check().Return(expectedDetails[1]).Times(1)

		overallStatus := healthService.ReadinessCheck(ctx)

		assert.Equal(t, health.StatusDegraded, overallStatus.OverallStatus) // Overall status should be DEGRADED
		assert.Len(t, overallStatus.Checks, 2)
		for _, expected := range expectedDetails {
			assert.Contains(t, overallStatus.Checks, expected)
		}
	})

	t.Run("Mix of DOWN and DEGRADED (DOWN takes precedence)", func(t *testing.T) {
		mockCheckers, healthService := setup(ctrl, 3)
		errStrDown := "critical service down"
		errStrDegraded := "minor service degraded"
		expectedDetails := []health.CheckDetail{
			{Name: "CriticalService", Status: health.StatusDown, Error: &errStrDown},      // DOWN
			{Name: "MinorService", Status: health.StatusDegraded, Error: &errStrDegraded}, // DEGRADED
			{Name: "AnotherService", Status: health.StatusUp, Error: nil},                 // UP
		}

		mockCheckers[0].EXPECT().Check().Return(expectedDetails[0]).Times(1)
		mockCheckers[1].EXPECT().Check().Return(expectedDetails[1]).Times(1)
		mockCheckers[2].EXPECT().Check().Return(expectedDetails[2]).Times(1)

		overallStatus := healthService.ReadinessCheck(ctx)

		assert.Equal(t, health.StatusDown, overallStatus.OverallStatus) // Overall status should be DOWN
		assert.Len(t, overallStatus.Checks, 3)
		for _, expected := range expectedDetails {
			assert.Contains(t, overallStatus.Checks, expected)
		}
	})

	t.Run("No checkers configured", func(t *testing.T) {
		// No checkers are provided to the service
		_, healthService := setup(ctrl, 0)

		overallStatus := healthService.ReadinessCheck(ctx)

		assert.Equal(t, health.StatusUp, overallStatus.OverallStatus) // Default to UP if no checks are performed
		assert.Empty(t, overallStatus.Checks)                         // No check details
		assert.WithinDuration(t, time.Now(), overallStatus.Timestamp, 5*time.Second)
	})

	t.Run("Checker returns an error message", func(t *testing.T) {
		mockCheckers, healthService := setup(ctrl, 1)
		errMessage := "simulated error during check"
		expectedDetail := health.CheckDetail{
			Name:   "FaultyChecker",
			Status: health.StatusDown, // Can be DOWN or DEGRADED with an error message
			Error:  &errMessage,
		}

		mockCheckers[0].EXPECT().Check().Return(expectedDetail).Times(1)

		overallStatus := healthService.ReadinessCheck(ctx)
		assert.Equal(t, health.StatusDown, overallStatus.OverallStatus)
		assert.Len(t, overallStatus.Checks, 1)
		assert.Equal(t, expectedDetail.Name, overallStatus.Checks[0].Name)
		assert.Equal(t, expectedDetail.Status, overallStatus.Checks[0].Status)
		assert.Equal(t, *expectedDetail.Error, *overallStatus.Checks[0].Error)
	})
}
