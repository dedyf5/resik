// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"github.com/dedyf5/resik/ctx/app"
	"github.com/dedyf5/resik/utils/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	log  *log.Log
	path string
}

type Obj struct {
	AppName string
	URI     string
}

func NewLog(log *log.Log, uri string) *Log {
	return &Log{
		log:  log,
		path: uri,
	}
}

func (l *Log) Error(msg string) {
	obj := &Obj{
		AppName: app.NameHTTP.String(),
		URI:     l.path,
	}
	l.log.Logger.Error(msg, zap.Inline(obj))
}

func (l *Log) Warn(msg string) {
	obj := &Obj{
		AppName: app.NameHTTP.String(),
		URI:     l.path,
	}
	l.log.Logger.Warn(msg, zap.Inline(obj))
}

func (l *Log) Info(msg string) {
	obj := &Obj{
		AppName: app.NameHTTP.String(),
		URI:     l.path,
	}
	l.log.Logger.Info(msg, zap.Inline(obj))
}

func (l *Log) Debug(msg string) {
	obj := &Obj{
		AppName: app.NameHTTP.String(),
		URI:     l.path,
	}
	l.log.Logger.Debug(msg, zap.Inline(obj))
}

func (o *Obj) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("app", o.AppName)
	enc.AddString("uri", o.URI)
	return nil
}
