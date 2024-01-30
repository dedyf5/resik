// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"github.com/dedyf5/resik/ctx/app"
	logUtil "github.com/dedyf5/resik/utils/log"
)

type Log struct {
	log  *logUtil.Log
	path string
}

func NewLog(log *logUtil.Log, uri string) *Log {
	return &Log{
		log:  log,
		path: uri,
	}
}

func (l *Log) Error(msg string) {
	l.log.Error(&logUtil.Service{
		AppName: app.NameHTTP.String(),
		Path:    l.path,
	}, msg)
}

func (l *Log) Warn(msg string) {
	l.log.Warn(&logUtil.Service{
		AppName: app.NameHTTP.String(),
		Path:    l.path,
	}, msg)
}

func (l *Log) Info(msg string) {
	l.log.Info(&logUtil.Service{
		AppName: app.NameHTTP.String(),
		Path:    l.path,
	}, msg)
}

func (l *Log) Debug(msg string) {
	l.log.Debug(&logUtil.Service{
		AppName: app.NameHTTP.String(),
		Path:    l.path,
	}, msg)
}
