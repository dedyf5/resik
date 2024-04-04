package datetime

import (
	"net/http"
	"time"

	resPkg "github.com/dedyf5/resik/pkg/response"
)

type Format string

const (
	FormatyyyyMMddHHmmss Format = "2006-01-02 15:04:05"
)

func (f *Format) ToString() string {
	return string(*f)
}

func FromString(val string, format Format) (res *time.Time, status *resPkg.Status) {
	datetime, err := time.Parse(format.ToString(), val)
	if err != nil {
		return nil, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	localTime := LocalTime(datetime)
	return &localTime, nil
}

func LocalTime(from time.Time) time.Time {
	_, offset := time.Now().Zone()
	diff := time.Duration(offset) * time.Second
	return from.Add(-diff)
}
