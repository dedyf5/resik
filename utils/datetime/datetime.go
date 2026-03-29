package datetime

import (
	"net/http"
	"time"

	"github.com/dedyf5/resik/ctx"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func FromString(val string, format string, c *ctx.Ctx) (res *time.Time, err *resPkg.Status) {
	datetime, errParse := time.Parse(format, val)
	if errParse != nil {
		return nil, resPkg.NewStatusMessage(
			http.StatusInternalServerError,
			c.Lang().GetByMessageID("invalid_time_format"),
			errParse,
		)
	}

	return &datetime, nil
}
