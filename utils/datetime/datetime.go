package datetime

import (
	"net/http"
	"time"

	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func FromString(val string, format string, c *ctx.Ctx) (res *time.Time, err *resPkg.Status) {
	datetime, errParse := time.Parse(format, val)
	if errParse != nil {
		return nil, resPkg.NewStatusMessage(
			http.StatusInternalServerError,
			term.InvalidTimeFormat.Localize(c.Lang().Localizer),
			errParse,
		)
	}

	return &datetime, nil
}
