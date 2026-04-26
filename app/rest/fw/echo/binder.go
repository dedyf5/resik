// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/lang/term"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v5"
)

type Binder interface {
	Bind(c *echo.Context, i any) error
	BindHeaders(c *echo.Context, i any) error
	BindBody(c *echo.Context, i any) error
	BindQueryParams(c *echo.Context, i any) error
	BindPathParams(c *echo.Context, i any) error
	ParamValidator(c *echo.Context, i any) error
	JSONErrorFormatter(c *echo.Context, err error) error
}

var binder *bind
var once sync.Once

type bind struct {
	def *echo.DefaultBinder
}

func NewBinder() Binder {
	once.Do(func() {
		binder = &bind{&echo.DefaultBinder{}}
	})
	return binder
}

func (b *bind) BindHeaders(c *echo.Context, i any) error {
	return echo.BindHeaders(c, i)
}

func (b *bind) BindBody(c *echo.Context, i any) error {
	return echo.BindBody(c, i)
}

func (b *bind) BindQueryParams(c *echo.Context, i any) error {
	return echo.BindQueryParams(c, i)
}

func (b *bind) BindPathParams(c *echo.Context, i any) error {
	return echo.BindPathValues(c, i)
}

func (b *bind) Bind(c *echo.Context, i any) error {
	if err := b.ParamValidator(c, i); err != nil {
		return err
	}
	if err := b.def.Bind(c, i); err != nil {
		return b.JSONErrorFormatter(c, err)
	}
	return nil
}

func (b *bind) ParamValidator(c *echo.Context, i any) error {
	if i == nil {
		return nil
	}

	if langString := c.Request().URL.Query().Get(langCtx.ContextKey.String()); langString != "" {
		if _, err := langCtx.LanguageIsAvailable(langString); err != nil {
			return err
		}
	}

	lang, err := langCtx.FromContext(c.Request().Context())
	if err != nil {
		return err
	}

	df := reflect.TypeOf(i)
	if df.Kind() == reflect.Pointer {
		df = df.Elem()
	}

	if df.Kind() != reflect.Struct {
		return nil
	}

	validate := func(tn string, getter func(string) string) error {
		for f := range df.Fields() {
			ftq := f.Tag.Get(tn)
			if ftq == "" {
				continue
			}

			fnames := strings.Split(ftq, ",")
			if len(fnames) == 0 {
				continue
			}

			fn := fnames[0]
			fv := getter(fn)
			if fv == "" {
				continue
			}

			ft := f.Type.Kind()
			if ft == reflect.Pointer {
				ft = f.Type.Elem().Kind()
			}

			if ft == reflect.String {
				continue
			}

			wn := map[reflect.Kind]bool{
				reflect.Int:    true,
				reflect.Int8:   true,
				reflect.Int16:  true,
				reflect.Int32:  true,
				reflect.Int64:  true,
				reflect.Uint:   true,
				reflect.Uint8:  true,
				reflect.Uint16: true,
				reflect.Uint32: true,
				reflect.Uint64: true,
			}
			if _, ok := wn[ft]; ok && regexp.MustCompile(`^-?\d+$`).Match([]byte(fv)) {
				continue
			}

			fln := map[reflect.Kind]bool{reflect.Float32: true, reflect.Float64: true}
			if _, ok := fln[ft]; ok && regexp.MustCompile(`^-?\d+\.?\d*$`).Match([]byte(fv)) {
				continue
			}

			if ft == reflect.Bool && regexp.MustCompile(`true|false|1|0$`).Match([]byte(fv)) {
				continue
			}

			nsc := map[reflect.Kind]bool{reflect.Slice: true, reflect.Array: true}
			if _, ok := nsc[ft]; ok {
				continue
			}

			message, technicalErr := errorTypeMessage(lang, fn, ft.String(), "string")

			return resPkg.NewStatusBadRequest(
				lang.LanguageReqOrDefault().String(),
				fn,
				message,
				"INVALID_TYPE",
				technicalErr,
			)
		}

		return nil
	}

	if err := validate("query", c.QueryParam); err != nil {
		return err
	}

	if err := validate("param", c.Param); err != nil {
		return err
	}

	if err := validate("form", c.FormValue); err != nil {
		return err
	}

	return nil
}

func (b *bind) JSONErrorFormatter(c *echo.Context, err error) error {
	lang, errLang := langCtx.FromContext(c.Request().Context())
	if errLang != nil {
		return errLang
	}
	regField := regexp.MustCompile(`field\=(.*?),`)
	fields := regField.FindStringSubmatch(err.Error())
	expectedReg := regexp.MustCompile(`expected\=(.*?),`)
	expecteds := expectedReg.FindStringSubmatch(err.Error())
	gotReg := regexp.MustCompile(`got\=(.*?),`)
	gots := gotReg.FindStringSubmatch(err.Error())
	if len(fields) > 1 && len(expecteds) > 1 && len(gots) > 1 {
		field := fields[1]
		message, technicalErr := errorTypeMessage(lang, field, expecteds[1], gots[1])
		return resPkg.NewStatusBadRequest(
			lang.LanguageReqOrDefault().String(),
			field,
			message,
			"INVALID_TYPE",
			technicalErr,
		)
	}

	return nil
}

func errorTypeMessage(lang *langCtx.Lang, field, expected, actual string) (message string, technicalErr error) {
	technicalErr = fmt.Errorf("field=%s,expected=%s,got=%s", field, expected, actual)
	message = term.ValidationTypeMessage.Localize(
		lang.Localizer,
		lang.GetValidationFieldNameWithQuote(field),
		expected,
		actual,
	)
	return
}
