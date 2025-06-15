// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"context"
	"sync"
	"time"

	"github.com/dedyf5/resik/core/health"
)

func (s *Service) LivenessCheck(c context.Context) (isLive bool, statusMessage string) {
	return true, "SERVING"
}

func (s *Service) ReadinessCheck(c context.Context) health.OverallHealthStatus {
	overallStatus := health.StatusUp
	var checkDetails []health.CheckDetail
	var wg sync.WaitGroup
	resultsChan := make(chan health.CheckDetail, len(s.checkers))

	for _, checker := range s.checkers {
		wg.Add(1)
		go func(chk health.Checker) {
			defer wg.Done()
			resultsChan <- chk.Check()
		}(checker)
	}

	wg.Wait()
	close(resultsChan)

	for detail := range resultsChan {
		checkDetails = append(checkDetails, detail)
		if detail.Status == health.StatusDown {
			overallStatus = health.StatusDown
		} else if detail.Status == health.StatusDegraded && overallStatus != health.StatusDown {
			overallStatus = health.StatusDegraded
		}
	}

	return health.OverallHealthStatus{
		OverallStatus: overallStatus,
		Timestamp:     time.Now(),
		Checks:        checkDetails,
	}
}
