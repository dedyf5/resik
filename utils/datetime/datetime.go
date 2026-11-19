// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package datetime provides utility functions for parsing, formatting,
// and manipulating date and time values within the application context.
package datetime

import (
	"net/http"
	"time"

	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

// FromString parses a formatted string into a *time.Time pointer using the specified layout format.
//
// If val is en empty string, it returns (nil, nil) without error.
// If parsing fails, it returns a localized error status wrapped in *resPkg.Status
func FromString(val string, format string, c *ctx.Ctx) (res *time.Time, err *resPkg.Status) {
	if val == "" {
		return nil, nil
	}

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
