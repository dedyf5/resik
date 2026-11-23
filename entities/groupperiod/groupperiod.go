// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package groupperiod defines time-grouping modes and period range helpers for
// grouping metrics and formatting date expressions across SQL backends.
package groupperiod

import (
	"time"
)

// Mode represents the granularity used to group time-series or period data.
type Mode string

const (
	ModeDay   Mode = "day"   // ModeDay groups period data by calendar day.
	ModeMonth Mode = "month" // ModeMonth groups period data by calendar month.
	ModeYear  Mode = "year"  // ModeYear groups period data by calendar year.
)

// GroupPeriod defines the grouping mode, start and end timestamps, and target timezone	for a query period.
type GroupPeriod struct {
	Mode          Mode
	DatetimeStart *time.Time
	DatetimeEnd   *time.Time
	Timezone      string
}

// DateFormatMySQL returns the MySQL DATE_FORMAT pattern corresponding to the Mode,
// or an empty string if the mode is unsupported.
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

// DateFormatPSQL returns the PostgreSQL TO_CHAR date format pattern corresponding to the Mode,
// or an empty string if the mode is unsupported.
func (m Mode) DateFormatPSQL() string {
	switch m {
	case ModeDay:
		return "YYYY-MM-DD"
	case ModeMonth:
		return "YYYY-MM"
	case ModeYear:
		return "YYYY"
	}
	return ""
}

// DatetimeStartString formats DatetimeStart in UTC as "YYYY-MM-DD HH:MM:SS".
func (g *GroupPeriod) DatetimeStartString() string {
	return g.DatetimeStart.UTC().Format("2006-01-02 15:04:05")
}

// DatetimeEndString formats DatetimeEnd in UTC as "YYYY-MM-DD HH:MM:SS".
func (g *GroupPeriod) DatetimeEndString() string {
	return g.DatetimeEnd.UTC().Format("2006-01-02 15:04:05")
}
