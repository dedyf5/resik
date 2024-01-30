// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package app

import (
	"github.com/dedyf5/resik/ctx/status"
	logUtil "github.com/dedyf5/resik/utils/log"
)

type Name string

const (
	NameHTTP Name = "http"
)

type IApp interface {
	Name() Name
	Location() string
	Status() status.IStatus
	Logger() logUtil.ILog
}

func (n Name) String() string {
	return string(n)
}
