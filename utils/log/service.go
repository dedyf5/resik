// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"go.uber.org/zap/zapcore"
)

type Service struct {
	AppName string
	Path    string
}

func (o *Service) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("app", o.AppName)
	enc.AddString("path", o.Path)
	return nil
}
