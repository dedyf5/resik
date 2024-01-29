// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package app

type Name string

const (
	NameHTTP Name = "http"
)

type IApp interface {
	Name() Name
	Location() string
	Status() IStatus
	Logger() ILog
}

type IStatus interface {
	IsError() bool
	Error() string
	MessageOrDefault() string
}

type ILog interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

func (n Name) String() string {
	return string(n)
}
