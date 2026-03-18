package datetime

import (
	"net/http"
	"time"

	resPkg "github.com/dedyf5/resik/pkg/response"
)

func FromString(val string, format string) (res *time.Time, status *resPkg.Status) {
	datetime, err := time.Parse(format, val)
	if err != nil {
		return nil, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}

	return &datetime, nil
}
