// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package app

import "github.com/dedyf5/resik/ctx/status"

type Name string

const (
	NameHTTP Name = "http"
)

type IApp interface {
	Name() Name
	Location() string
	Status() status.IStatus
	Logger() ILog
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
