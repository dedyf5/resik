package datetime

import (
	"net/http"
	"time"

	"github.com/dedyf5/resik/ctx"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func FromString(val string, format string, c *ctx.Ctx) (res *time.Time, status *resPkg.Status) {
	datetime, err := time.Parse(format, val)
	if err != nil {
		return nil, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			Message:    c.Lang().GetByMessageID("invalid_time_format"),
			CauseError: err,
		}
	}

	return &datetime, nil
}
