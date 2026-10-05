// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dedyf5/resik/core/health"
	checkEntity "github.com/dedyf5/resik/entities/check"
	repo "github.com/dedyf5/resik/repositories"
	mockRepo "github.com/dedyf5/resik/repositories/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// setup initializes gomock controller, mock checkers, and the health service.
// It returns the slice of mock checkers and the health service instance.
func setup(ctrl *gomock.Controller, numCheckers int) ([]*mockRepo.MockICheck, health.IService) {
	mockCheckers := make([]*mockRepo.MockICheck, numCheckers)
	healthCheckers := make([]repo.ICheck, numCheckers)
	for i := range numCheckers {
		mockCheckers[i] = mockRepo.NewMockICheck(ctrl)
		healthCheckers[i] = mockCheckers[i]
	}
	healthService := New(healthCheckers)
	return mockCheckers, healthService
}

func TestLivenessCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, healthService := setup(ctrl, 0)

	t.Run("TestLivenessCheck ALL-SUCCESS", func(t *testing.T) {
		isLive, statusMessage := healthService.LivenessCheck(context.Background())
		assert.True(t, isLive)
		assert.Equal(t, "SERVING", statusMessage)
	})
}

func TestReadinessCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	t.Run("TestReadinessCheck ALL-UP", func(t *testing.T) {
		mockCheckers, healthService := setup(ctrl, 3)
		expectedDetails := []checkEntity.CheckDetail{
			{Name: "Checker1", Status: checkEntity.StatusUp, Error: nil},
			{Name: "Checker2", Status: checkEntity.StatusUp, Error: nil},
			{Name: "Checker3", Status: checkEntity.StatusUp, Error: nil},
		}

		// Set expectations for each mock checker
		mockCheckers[0].EXPECT().Check().Return(expectedDetails[0]).Times(1)
		mockCheckers[1].EXPECT().Check().Return(expectedDetails[1]).Times(1)
		mockCheckers[2].EXPECT().Check().Return(expectedDetails[2]).Times(1)

		overallStatus := healthService.ReadinessCheck(ctx)

		assert.Equal(t, checkEntity.StatusUp, overallStatus.OverallStatus)
		assert.Len(t, overallStatus.Checks, 3)
		// Assert individual check details. Order is not guaranteed due to goroutines.
		// We use assert.Contains to check if each expected detail is present.
		for _, expected := range expectedDetails {
			assert.Contains(t, overallStatus.Checks, expected)
		}
		// Check if the timestamp is recent (within a few seconds)
		assert.WithinDuration(t, time.Now(), overallStatus.Timestamp, 5*time.Second)
	})

	t.Run("TestReadinessCheck ONE-DOWN", func(t *testing.T) {
		mockCheckers, healthService := setup(ctrl, 2)
		errStr := errors.New("database connection failed")
		expectedDetails := []checkEntity.CheckDetail{
			{Name: "DBChecker", Status: checkEntity.StatusDown, Error: errStr},
			{Name: "APIChecker", Status: checkEntity.StatusUp, Error: nil},
		}

		mockCheckers[0].EXPECT().Check().Return(expectedDetails[0]).Times(1)
		mockCheckers[1].EXPECT().Check().Return(expectedDetails[1]).Times(1)

		overallStatus := healthService.ReadinessCheck(ctx)

		assert.Equal(t, checkEntity.StatusDown, overallStatus.OverallStatus)
		assert.Len(t, overallStatus.Checks, 2)
		for _, expected := range expectedDetails {
			assert.Contains(t, overallStatus.Checks, expected)
		}
	})

	t.Run("TestReadinessCheck ONE-DEGRADED OTHERS-UP", func(t *testing.T) {
		mockCheckers, healthService := setup(ctrl, 2)
		errStr := errors.New("cache performance degraded")
		expectedDetails := []checkEntity.CheckDetail{
			{Name: "CacheChecker", Status: checkEntity.StatusDegraded, Error: errStr},
			{Name: "QueueChecker", Status: checkEntity.StatusUp, Error: nil},
		}

		mockCheckers[0].EXPECT().Check().Return(expectedDetails[0]).Times(1)
		mockCheckers[1].EXPECT().Check().Return(expectedDetails[1]).Times(1)

		overallStatus := healthService.ReadinessCheck(ctx)

		assert.Equal(t, checkEntity.StatusDegraded, overallStatus.OverallStatus)
		assert.Len(t, overallStatus.Checks, 2)
		for _, expected := range expectedDetails {
			assert.Contains(t, overallStatus.Checks, expected)
		}
	})

	t.Run("TestReadinessCheck MIXED-DOWN-DEGRADED-UP", func(t *testing.T) {
		mockCheckers, healthService := setup(ctrl, 3)
		errStrDown := errors.New("critical service down")
		errStrDegraded := errors.New("minor service degraded")
		expectedDetails := []checkEntity.CheckDetail{
			{Name: "CriticalService", Status: checkEntity.StatusDown, Error: errStrDown},
			{Name: "MinorService", Status: checkEntity.StatusDegraded, Error: errStrDegraded},
			{Name: "AnotherService", Status: checkEntity.StatusUp, Error: nil},
		}

		mockCheckers[0].EXPECT().Check().Return(expectedDetails[0]).Times(1)
		mockCheckers[1].EXPECT().Check().Return(expectedDetails[1]).Times(1)
		mockCheckers[2].EXPECT().Check().Return(expectedDetails[2]).Times(1)

		overallStatus := healthService.ReadinessCheck(ctx)

		assert.Equal(t, checkEntity.StatusDown, overallStatus.OverallStatus)
		assert.Len(t, overallStatus.Checks, 3)
		for _, expected := range expectedDetails {
			assert.Contains(t, overallStatus.Checks, expected)
		}
	})

	t.Run("TestReadinessCheck NO-CHECKERS", func(t *testing.T) {
		_, healthService := setup(ctrl, 0)

		overallStatus := healthService.ReadinessCheck(ctx)

		assert.Equal(t, checkEntity.StatusUp, overallStatus.OverallStatus)
		assert.Empty(t, overallStatus.Checks)
		assert.WithinDuration(t, time.Now(), overallStatus.Timestamp, 5*time.Second)
	})

	t.Run("TestReadinessCheck ERROR", func(t *testing.T) {
		mockCheckers, healthService := setup(ctrl, 1)
		errMessage := errors.New("simulated error during check")
		expectedDetail := checkEntity.CheckDetail{
			Name:   "FaultyChecker",
			Status: checkEntity.StatusDown, // Can be DOWN or DEGRADED with an error message
			Error:  errMessage,
		}

		mockCheckers[0].EXPECT().Check().Return(expectedDetail).Times(1)

		overallStatus := healthService.ReadinessCheck(ctx)
		assert.Equal(t, checkEntity.StatusDown, overallStatus.OverallStatus)
		assert.Len(t, overallStatus.Checks, 1)
		assert.Equal(t, expectedDetail.Name, overallStatus.Checks[0].Name)
		assert.Equal(t, expectedDetail.Status, overallStatus.Checks[0].Status)
		assert.Equal(t, expectedDetail.Error, overallStatus.Checks[0].Error)
	})
}
