// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package app

import (
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/ctx/status"
)

type Name string

const (
	NameHTTP Name = "http"
)

type IApp interface {
	Name() Name
	Path() string
	Status() status.IStatus
	Logger() logCtx.ILog
}

func (n Name) String() string {
	return string(n)
}
