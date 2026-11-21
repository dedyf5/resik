// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package groupperiod

import (
	"testing"
	"time"
)

func TestMode_DateFormatMySQL(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
		want string
	}{
		{name: "day", mode: ModeDay, want: "%Y-%m-%d"},
		{name: "month", mode: ModeMonth, want: "%Y-%m"},
		{name: "year", mode: ModeYear, want: "%Y"},
		{name: "unsupported", mode: Mode("week"), want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.mode.DateFormatMySQL(); got != tc.want {
				t.Errorf("DateFormatMySQL() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMode_DateFormatPSQL(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
		want string
	}{
		{name: "day", mode: ModeDay, want: "YYYY-MM-DD"},
		{name: "month", mode: ModeMonth, want: "YYYY-MM"},
		{name: "year", mode: ModeYear, want: "YYYY"},
		{name: "unsupported", mode: Mode("week"), want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.mode.DateFormatPSQL(); got != tc.want {
				t.Errorf("DateFormatPSQL() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGroupPeriod_DateTimeStrings(t *testing.T) {
	loc := time.FixedZone("WIB", 7*3600)
	start := time.Date(2026, time.November, 21, 14, 30, 5, 0, loc) // 07:30:05 UTC
	end := time.Date(2026, time.November, 21, 23, 59, 59, 0, loc)  // 16:59:59 UTC

	gp := &GroupPeriod{
		Mode:          ModeDay,
		DatetimeStart: &start,
		DatetimeEnd:   &end,
		Timezone:      "Asia/Jakarta",
	}

	if got, want := gp.DatetimeStartString(), "2026-11-21 07:30:05"; got != want {
		t.Errorf("DatetimeStartString() = %q, want %q", got, want)
	}
	if got, want := gp.DatetimeEndString(), "2026-11-21 16:59:59"; got != want {
		t.Errorf("DatetimeEndString() = %q, want %q", got, want)
	}
}
