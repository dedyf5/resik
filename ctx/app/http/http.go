// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"github.com/dedyf5/resik/ctx/app"
	"github.com/dedyf5/resik/ctx/status"
	logUtil "github.com/dedyf5/resik/utils/log"
)

type HTTP struct {
	log    *Log
	uri    string
	status status.Status
}

func NewHTTP(log *logUtil.Log, uri string) *HTTP {
	return &HTTP{
		uri: uri,
		log: NewLog(log, uri),
	}
}

func (h *HTTP) Name() app.Name {
	return app.NameHTTP
}

func (h *HTTP) Location() string {
	return h.uri
}

func (h *HTTP) Status() status.IStatus {
	return &h.status
}

func (h *HTTP) Logger() logUtil.ILog {
	return h.log
}
