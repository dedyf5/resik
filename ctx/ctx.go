// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package ctx

import (
	"context"

	"github.com/dedyf5/resik/ctx/app"
	httpApp "github.com/dedyf5/resik/ctx/app/http"
	lang "github.com/dedyf5/resik/ctx/lang"
	logUtil "github.com/dedyf5/resik/utils/log"
)

type Ctx struct {
	App     app.IApp
	Lang    *lang.Lang
	Context context.Context
}

func NewHTTP(ctx context.Context, log *logUtil.Log, lang *lang.Lang, uri string) *Ctx {
	return &Ctx{
		App:     httpApp.NewHTTP(log, uri),
		Context: ctx,
		Lang:    lang,
	}
}

func (c *Ctx) IsError() bool {
	return c.App.Status().IsError()
}

func (c *Ctx) Error() string {
	return c.App.Status().Error()
}

func (c *Ctx) MessageOrDefault() string {
	return c.App.Status().MessageOrDefault()
}
