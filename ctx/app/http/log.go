// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"github.com/dedyf5/resik/ctx/app"
	logCtx "github.com/dedyf5/resik/ctx/log"
)

type Log struct {
	log  *logCtx.Log
	path string
}

func NewLog(log *logCtx.Log, uri string) *Log {
	return &Log{
		log:  log,
		path: uri,
	}
}

func (l *Log) Error(msg string) {
	l.log.Error(&logCtx.Service{
		AppName: app.NameHTTP.String(),
		Path:    l.path,
	}, msg)
}

func (l *Log) Warn(msg string) {
	l.log.Warn(&logCtx.Service{
		AppName: app.NameHTTP.String(),
		Path:    l.path,
	}, msg)
}

func (l *Log) Info(msg string) {
	l.log.Info(&logCtx.Service{
		AppName: app.NameHTTP.String(),
		Path:    l.path,
	}, msg)
}

func (l *Log) Debug(msg string) {
	l.log.Debug(&logCtx.Service{
		AppName: app.NameHTTP.String(),
		Path:    l.path,
	}, msg)
}
