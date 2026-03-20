// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package groupperiod

import (
	"time"
)

type Mode string

const (
	ModeDay   Mode = "day"
	ModeMonth Mode = "month"
	ModeYear  Mode = "year"
)

type GroupPeriod struct {
	Mode          Mode
	DatetimeStart *time.Time
	DatetimeEnd   *time.Time
	Timezone      string
}

func (m Mode) DateFormatMySQL() string {
	switch m {
	case ModeDay:
		return "%Y-%m-%d"
	case ModeMonth:
		return "%Y-%m"
	case ModeYear:
		return "%Y"
	}
	return ""
}

func (g *GroupPeriod) DatetimeStartString() string {
	return g.DatetimeStart.UTC().Format("2006-01-02 15:04:05")
}

func (g *GroupPeriod) DatetimeEndString() string {
	return g.DatetimeEnd.UTC().Format("2006-01-02 15:04:05")
}
