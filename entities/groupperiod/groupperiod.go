// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package groupperiod

type Mode string

const (
	ModeDay   Mode = "day"
	ModeMonth Mode = "month"
	ModeYear  Mode = "year"
)

type GroupPeriod struct {
	Mode          Mode
	DatetimeStart string // yyyy-MM-dd HH:mm:ss
	DatetimeEnd   string // yyyy-MM-dd HH:mm:ss
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
