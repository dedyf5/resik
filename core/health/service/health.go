// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"context"
	"sync"
	"time"

	checkEntity "github.com/dedyf5/resik/entities/check"
	repo "github.com/dedyf5/resik/repositories"
)

func (s *Service) LivenessCheck(c context.Context) (isLive bool, statusMessage string) {
	return true, "SERVING"
}

func (s *Service) ReadinessCheck(c context.Context) *checkEntity.OverallHealthStatus {
	var wg sync.WaitGroup
	resultsChan := make(chan checkEntity.CheckDetail, len(s.checkers))

	for _, checker := range s.checkers {
		wg.Add(1)
		go func(chk repo.ICheck) {
			defer wg.Done()
			resultsChan <- chk.Check()
		}(checker)
	}

	wg.Wait()
	close(resultsChan)

	var checkDetails []checkEntity.CheckDetail
	for detail := range resultsChan {
		checkDetails = append(checkDetails, detail)
	}

	return checkEntity.NewOverallHealthStatus(time.Now(), checkDetails...)
}
