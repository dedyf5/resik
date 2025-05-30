// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/entities/common"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v4"
)

type Binder interface {
	Bind(i interface{}, c echo.Context) error
	BindHeaders(c echo.Context, i interface{}) error
	BindBody(c echo.Context, i interface{}) error
	BindQueryParams(c echo.Context, i interface{}) error
	BindPathParams(c echo.Context, i interface{}) error
	ParamValidator(c echo.Context, i interface{}) error
	JSONErrorFormatter(c echo.Context, err error) error
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

func (b *bind) BindHeaders(c echo.Context, i interface{}) error {
	return b.def.BindHeaders(c, i)
}

func (b *bind) BindBody(c echo.Context, i interface{}) error {
	return b.def.BindBody(c, i)
}

func (b *bind) BindQueryParams(c echo.Context, i interface{}) error {
	return b.def.BindQueryParams(c, i)
}

func (b *bind) BindPathParams(c echo.Context, i interface{}) error {
	return b.def.BindPathParams(c, i)
}

func (b *bind) Bind(i interface{}, c echo.Context) error {
	if err := b.ParamValidator(c, i); err != nil {
		return err
	}
	if err := b.def.Bind(i, c); err != nil {
		return b.JSONErrorFormatter(c, err)
	}
	return nil
}

func (b *bind) ParamValidator(c echo.Context, i interface{}) error {
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
	if df.Kind() == reflect.Ptr {
		df = df.Elem()
	}

	if df.Kind() != reflect.Struct {
		return nil
	}

	validate := func(tn string, getter func(string) string) error {
		for i := 0; i < df.NumField(); i++ {
			f := df.Field(i)
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

			ft := df.Field(i).Type.Kind()
			if ft == reflect.Ptr {
				ft = df.Field(i).Type.Elem().Kind()
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

			msg := lang.GetByTemplateData("validation_type_message", common.Map{
				"field":    fn,
				"expected": ft,
				"actual":   "string",
			})
			return &resPkg.Status{
				Code:    http.StatusBadRequest,
				Message: msg,
				Detail: map[string]string{
					fn: msg,
				},
			}
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

func (b *bind) JSONErrorFormatter(c echo.Context, err error) error {
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
		errMap := map[string]string{}
		field := fields[1]
		errMap[field] = lang.GetByTemplateData("validation_type_message", common.Map{
			"field":    field,
			"expected": expecteds[1],
			"actual":   gots[1],
		})
		return &resPkg.Status{
			Code:    http.StatusBadRequest,
			Message: errMap[field],
			Detail:  errMap,
		}
	}

	return nil
}
