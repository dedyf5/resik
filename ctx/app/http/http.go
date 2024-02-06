// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"github.com/dedyf5/resik/ctx/app"
	logCtx "github.com/dedyf5/resik/ctx/log"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

type HTTP struct {
	log    *Log
	uri    string
	status resPkg.Status
}

func NewHTTP(log *logCtx.Log, uri string) *HTTP {
	return &HTTP{
		uri: uri,
		log: NewLog(log, uri),
	}
}

func (h *HTTP) Name() app.Name {
	return app.NameHTTP
}

func (h *HTTP) Path() string {
	return h.uri
}

func (h *HTTP) Status() resPkg.IStatus {
	return &h.status
}

func (h *HTTP) Logger() logCtx.ILog {
	return h.log
}
