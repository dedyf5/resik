// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package translations

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// reference: https://github.com/go-playground/validator/blob/master/translations/id/id.go
var Indonesian []translation = []translation{
	{
		Tag:         "required",
		Translation: "{0} wajib diisi",
		Override:    false,
	},
	{
		Tag: "len",
		CustomRegisFunc: func(ut ut.Translator) (err error) {

			if err = ut.Add("len-string", "panjang {0} harus {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("len-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("len-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("len-number", "{0} harus sama dengan {1}", false); err != nil {
				return
			}

			if err = ut.Add("len-items", "{0} harus berisi {1}", false); err != nil {
				return
			}
			// if err = ut.AddCardinal("len-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("len-items-item", "{0} item", locales.PluralRuleOther, false); err != nil {
				return
			}

			return

		},
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			var err error
			var t string

			var digits uint64
			var kind reflect.Kind

			if idx := strings.Index(fe.Param(), "."); idx != -1 {
				digits = uint64(len(fe.Param()[idx+1:]))
			}

			f64, err := strconv.ParseFloat(fe.Param(), 64)
			if err != nil {
				goto END
			}

			kind = fe.Kind()
			if kind == reflect.Ptr {
				kind = fe.Type().Elem().Kind()
			}

			switch kind {
			case reflect.String:

				var c string

				c, err = ut.C("len-string-character", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("len-string", fe.Field(), c)

			case reflect.Slice, reflect.Map, reflect.Array:
				var c string

				c, err = ut.C("len-items-item", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("len-items", fe.Field(), c)

			default:
				t, err = ut.T("len-number", fe.Field(), ut.FmtNumber(f64, digits))
			}

		END:
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %s", err)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag: "min",
		CustomRegisFunc: func(ut ut.Translator) (err error) {

			if err = ut.Add("min-string", "panjang minimal {0} adalah {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("min-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("min-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("min-number", "{0} harus {1} atau lebih besar", false); err != nil {
				return
			}

			if err = ut.Add("min-items", "panjang minimal {0} adalah {1}", false); err != nil {
				return
			}
			// if err = ut.AddCardinal("min-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("min-items-item", "{0} item", locales.PluralRuleOther, false); err != nil {
				return
			}

			return

		},
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			var err error
			var t string

			var digits uint64
			var kind reflect.Kind

			if idx := strings.Index(fe.Param(), "."); idx != -1 {
				digits = uint64(len(fe.Param()[idx+1:]))
			}

			f64, err := strconv.ParseFloat(fe.Param(), 64)
			if err != nil {
				goto END
			}

			kind = fe.Kind()
			if kind == reflect.Ptr {
				kind = fe.Type().Elem().Kind()
			}

			switch kind {
			case reflect.String:

				var c string

				c, err = ut.C("min-string-character", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("min-string", fe.Field(), c)

			case reflect.Slice, reflect.Map, reflect.Array:
				var c string

				c, err = ut.C("min-items-item", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("min-items", fe.Field(), c)

			default:
				t, err = ut.T("min-number", fe.Field(), ut.FmtNumber(f64, digits))
			}

		END:
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %s", err)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag: "max",
		CustomRegisFunc: func(ut ut.Translator) (err error) {

			if err = ut.Add("max-string", "panjang maksimal {0} adalah {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("max-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("max-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("max-number", "{0} harus {1} atau kurang", false); err != nil {
				return
			}

			if err = ut.Add("max-items", "{0} harus berisi maksimal {1}", false); err != nil {
				return
			}
			// if err = ut.AddCardinal("max-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("max-items-item", "{0} item", locales.PluralRuleOther, false); err != nil {
				return
			}

			return

		},
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			var err error
			var t string

			var digits uint64
			var kind reflect.Kind

			if idx := strings.Index(fe.Param(), "."); idx != -1 {
				digits = uint64(len(fe.Param()[idx+1:]))
			}

			f64, err := strconv.ParseFloat(fe.Param(), 64)
			if err != nil {
				goto END
			}

			kind = fe.Kind()
			if kind == reflect.Ptr {
				kind = fe.Type().Elem().Kind()
			}

			switch kind {
			case reflect.String:

				var c string

				c, err = ut.C("max-string-character", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("max-string", fe.Field(), c)

			case reflect.Slice, reflect.Map, reflect.Array:
				var c string

				c, err = ut.C("max-items-item", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("max-items", fe.Field(), c)

			default:
				t, err = ut.T("max-number", fe.Field(), ut.FmtNumber(f64, digits))
			}

		END:
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %s", err)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "eq",
		Translation: "{0} tidak sama dengan {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "ne",
		Translation: "{0} tidak sama dengan {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag: "lt",
		CustomRegisFunc: func(ut ut.Translator) (err error) {

			if err = ut.Add("lt-string", "panjang {0} harus kurang dari {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("lt-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("lt-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lt-number", "{0} harus kurang dari {1}", false); err != nil {
				return
			}

			if err = ut.Add("lt-items", "{0} harus berisi kurang dari {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("lt-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("lt-items-item", "{0} item", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lt-datetime", "{0} harus kurang dari tanggal & waktu saat ini", false); err != nil {
				return
			}

			return

		},
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			var err error
			var t string
			var f64 float64
			var digits uint64
			var kind reflect.Kind

			fn := func() (err error) {

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err = strconv.ParseFloat(fe.Param(), 64)

				return
			}

			kind = fe.Kind()
			if kind == reflect.Ptr {
				kind = fe.Type().Elem().Kind()
			}

			switch kind {
			case reflect.String:

				var c string

				err = fn()
				if err != nil {
					goto END
				}

				c, err = ut.C("lt-string-character", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("lt-string", fe.Field(), c)

			case reflect.Slice, reflect.Map, reflect.Array:
				var c string

				err = fn()
				if err != nil {
					goto END
				}

				c, err = ut.C("lt-items-item", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("lt-items", fe.Field(), c)

			case reflect.Struct:
				if fe.Type() != reflect.TypeOf(time.Time{}) {
					err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
					goto END
				}

				t, err = ut.T("lt-datetime", fe.Field())

			default:
				err = fn()
				if err != nil {
					goto END
				}

				t, err = ut.T("lt-number", fe.Field(), ut.FmtNumber(f64, digits))
			}

		END:
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %s", err)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag: "lte",
		CustomRegisFunc: func(ut ut.Translator) (err error) {

			if err = ut.Add("lte-string", "panjang maksimal {0} adalah {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("lte-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("lte-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lte-number", "{0} harus {1} atau kurang", false); err != nil {
				return
			}

			if err = ut.Add("lte-items", "{0} harus berisi maksimal {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("lte-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("lte-items-item", "{0} item", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lte-datetime", "{0} harus kurang dari atau sama dengan tanggal & waktu saat ini", false); err != nil {
				return
			}

			return
		},
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			var err error
			var t string
			var f64 float64
			var digits uint64
			var kind reflect.Kind

			fn := func() (err error) {

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err = strconv.ParseFloat(fe.Param(), 64)

				return
			}

			kind = fe.Kind()
			if kind == reflect.Ptr {
				kind = fe.Type().Elem().Kind()
			}

			switch kind {
			case reflect.String:

				var c string

				err = fn()
				if err != nil {
					goto END
				}

				c, err = ut.C("lte-string-character", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("lte-string", fe.Field(), c)

			case reflect.Slice, reflect.Map, reflect.Array:
				var c string

				err = fn()
				if err != nil {
					goto END
				}

				c, err = ut.C("lte-items-item", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("lte-items", fe.Field(), c)

			case reflect.Struct:
				if fe.Type() != reflect.TypeOf(time.Time{}) {
					err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
					goto END
				}

				t, err = ut.T("lte-datetime", fe.Field())

			default:
				err = fn()
				if err != nil {
					goto END
				}

				t, err = ut.T("lte-number", fe.Field(), ut.FmtNumber(f64, digits))
			}

		END:
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %s", err)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag: "gt",
		CustomRegisFunc: func(ut ut.Translator) (err error) {

			if err = ut.Add("gt-string", "panjang {0} harus lebih dari {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("gt-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("gt-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gt-number", "{0} harus lebih besar dari {1}", false); err != nil {
				return
			}

			if err = ut.Add("gt-items", "{0} harus berisi lebih dari {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("gt-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("gt-items-item", "{0} item", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gt-datetime", "{0} harus lebih besar dari tanggal & waktu saat ini", false); err != nil {
				return
			}

			return
		},
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			var err error
			var t string
			var f64 float64
			var digits uint64
			var kind reflect.Kind

			fn := func() (err error) {

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err = strconv.ParseFloat(fe.Param(), 64)

				return
			}

			kind = fe.Kind()
			if kind == reflect.Ptr {
				kind = fe.Type().Elem().Kind()
			}

			switch kind {
			case reflect.String:

				var c string

				err = fn()
				if err != nil {
					goto END
				}

				c, err = ut.C("gt-string-character", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("gt-string", fe.Field(), c)

			case reflect.Slice, reflect.Map, reflect.Array:
				var c string

				err = fn()
				if err != nil {
					goto END
				}

				c, err = ut.C("gt-items-item", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("gt-items", fe.Field(), c)

			case reflect.Struct:
				if fe.Type() != reflect.TypeOf(time.Time{}) {
					err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
					goto END
				}

				t, err = ut.T("gt-datetime", fe.Field())

			default:
				err = fn()
				if err != nil {
					goto END
				}

				t, err = ut.T("gt-number", fe.Field(), ut.FmtNumber(f64, digits))
			}

		END:
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %s", err)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag: "gte",
		CustomRegisFunc: func(ut ut.Translator) (err error) {

			if err = ut.Add("gte-string", "panjang minimal {0} adalah {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("gte-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("gte-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gte-number", "{0} harus {1} atau lebih besar", false); err != nil {
				return
			}

			if err = ut.Add("gte-items", "{0} harus berisi setidaknya {1}", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("gte-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("gte-items-item", "{0} item", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gte-datetime", "{0} harus lebih besar dari atau sama dengan tanggal & waktu saat ini", false); err != nil {
				return
			}

			return
		},
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			var err error
			var t string
			var f64 float64
			var digits uint64
			var kind reflect.Kind

			fn := func() (err error) {

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err = strconv.ParseFloat(fe.Param(), 64)

				return
			}

			kind = fe.Kind()
			if kind == reflect.Ptr {
				kind = fe.Type().Elem().Kind()
			}

			switch kind {
			case reflect.String:

				var c string

				err = fn()
				if err != nil {
					goto END
				}

				c, err = ut.C("gte-string-character", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("gte-string", fe.Field(), c)

			case reflect.Slice, reflect.Map, reflect.Array:
				var c string

				err = fn()
				if err != nil {
					goto END
				}

				c, err = ut.C("gte-items-item", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}

				t, err = ut.T("gte-items", fe.Field(), c)

			case reflect.Struct:
				if fe.Type() != reflect.TypeOf(time.Time{}) {
					err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
					goto END
				}

				t, err = ut.T("gte-datetime", fe.Field())

			default:
				err = fn()
				if err != nil {
					goto END
				}

				t, err = ut.T("gte-number", fe.Field(), ut.FmtNumber(f64, digits))
			}

		END:
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %s", err)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "eqfield",
		Translation: "{0} harus sama dengan {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "eqcsfield",
		Translation: "{0} harus sama dengan {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "necsfield",
		Translation: "{0} tidak sama dengan {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "gtcsfield",
		Translation: "{0} harus lebih besar dari {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "gtecsfield",
		Translation: "{0} harus lebih besar dari atau sama dengan {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "ltcsfield",
		Translation: "{0} harus kurang dari {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "ltecsfield",
		Translation: "{0} harus kurang dari atau sama dengan {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "nefield",
		Translation: "{0} tidak sama dengan {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "gtfield",
		Translation: "{0} harus lebih besar dari {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "gtefield",
		Translation: "{0} harus lebih besar dari atau sama dengan {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "ltfield",
		Translation: "{0} harus kurang dari {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "ltefield",
		Translation: "{0} harus kurang dari atau sama dengan {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "alpha",
		Translation: "{0} hanya dapat berisi karakter abjad",
		Override:    false,
	},
	{
		Tag:         "alphanum",
		Translation: "{0} hanya dapat berisi karakter alfanumerik",
		Override:    false,
	},
	{
		Tag:         "numeric",
		Translation: "{0} harus berupa nilai numerik yang valid",
		Override:    false,
	},
	{
		Tag:         "number",
		Translation: "{0} harus berupa angka yang valid",
		Override:    false,
	},
	{
		Tag:         "hexadecimal",
		Translation: "{0} harus berupa heksadesimal yang valid",
		Override:    false,
	},
	{
		Tag:         "hexcolor",
		Translation: "{0} harus berupa warna HEX yang valid",
		Override:    false,
	},
	{
		Tag:         "rgb",
		Translation: "{0} harus berupa warna RGB yang valid",
		Override:    false,
	},
	{
		Tag:         "rgba",
		Translation: "{0} harus berupa warna RGBA yang valid",
		Override:    false,
	},
	{
		Tag:         "hsl",
		Translation: "{0} harus berupa warna HSL yang valid",
		Override:    false,
	},
	{
		Tag:         "hsla",
		Translation: "{0} harus berupa warna HSLA yang valid",
		Override:    false,
	},
	{
		Tag:         "email",
		Translation: "{0} harus berupa alamat email yang valid",
		Override:    false,
	},
	{
		Tag:         "url",
		Translation: "{0} harus berupa URL yang valid",
		Override:    false,
	},
	{
		Tag:         "uri",
		Translation: "{0} harus berupa URI yang valid",
		Override:    false,
	},
	{
		Tag:         "base64",
		Translation: "{0} harus berupa string Base64 yang valid",
		Override:    false,
	},
	{
		Tag:         "contains",
		Translation: "{0} harus berisi teks '{1}'",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "containsany",
		Translation: "{0} harus berisi setidaknya salah satu karakter berikut '{1}'",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "excludes",
		Translation: "{0} tidak boleh berisi teks '{1}'",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "excludesall",
		Translation: "{0} tidak boleh berisi salah satu karakter berikut '{1}'",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "excludesrune",
		Translation: "{0} tidak boleh berisi '{1}'",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "isbn",
		Translation: "{0} harus berupa nomor ISBN yang valid",
		Override:    false,
	},
	{
		Tag:         "isbn10",
		Translation: "{0} harus berupa nomor ISBN-10 yang valid",
		Override:    false,
	},
	{
		Tag:         "isbn13",
		Translation: "{0} harus berupa nomor ISBN-13 yang valid",
		Override:    false,
	},
	{
		Tag:         "issn",
		Translation: "{0} harus berupa nomor ISSN yang valid",
		Override:    false,
	},
	{
		Tag:         "uuid",
		Translation: "{0} harus berupa UUID yang valid",
		Override:    false,
	},
	{
		Tag:         "uuid3",
		Translation: "{0} harus berupa UUID versi 3 yang valid",
		Override:    false,
	},
	{
		Tag:         "uuid4",
		Translation: "{0} harus berupa UUID versi 4 yang valid",
		Override:    false,
	},
	{
		Tag:         "uuid5",
		Translation: "{0} harus berupa UUID versi 5 yang valid",
		Override:    false,
	},
	{
		Tag:         "ulid",
		Translation: "{0} harus berupa ULID yang valid",
		Override:    false,
	},
	{
		Tag:         "ascii",
		Translation: "{0} hanya boleh berisi karakter ascii",
		Override:    false,
	},
	{
		Tag:         "printascii",
		Translation: "{0} hanya boleh berisi karakter ascii yang dapat dicetak",
		Override:    false,
	},
	{
		Tag:         "multibyte",
		Translation: "{0} harus berisi karakter multibyte",
		Override:    false,
	},
	{
		Tag:         "datauri",
		Translation: "{0} harus berisi URI Data yang valid",
		Override:    false,
	},
	{
		Tag:         "latitude",
		Translation: "{0} harus berisi koordinat lintang yang valid",
		Override:    false,
	},
	{
		Tag:         "longitude",
		Translation: "{0} harus berisi koordinat bujur yang valid",
		Override:    false,
	},
	{
		Tag:         "ssn",
		Translation: "{0} harus berupa nomor SSN yang valid",
		Override:    false,
	},
	{
		Tag:         "ipv4",
		Translation: "{0} harus berupa alamat IPv4 yang valid",
		Override:    false,
	},
	{
		Tag:         "ipv6",
		Translation: "{0} harus berupa alamat IPv6 yang valid",
		Override:    false,
	},
	{
		Tag:         "ip",
		Translation: "{0} harus berupa alamat IP yang valid",
		Override:    false,
	},
	{
		Tag:         "cidr",
		Translation: "{0} harus berisi notasi CIDR yang valid",
		Override:    false,
	},
	{
		Tag:         "cidrv4",
		Translation: "{0} harus berisi notasi CIDR yang valid untuk alamat IPv4",
		Override:    false,
	},
	{
		Tag:         "cidrv6",
		Translation: "{0} harus berisi notasi CIDR yang valid untuk alamat IPv6",
		Override:    false,
	},
	{
		Tag:         "tcp_addr",
		Translation: "{0} harus berupa alamat TCP yang valid",
		Override:    false,
	},
	{
		Tag:         "tcp4_addr",
		Translation: "{0} harus berupa alamat TCP IPv4 yang valid",
		Override:    false,
	},
	{
		Tag:         "tcp6_addr",
		Translation: "{0} harus berupa alamat TCP IPv6 yang valid",
		Override:    false,
	},
	{
		Tag:         "udp_addr",
		Translation: "{0} harus berupa alamat UDP yang valid",
		Override:    false,
	},
	{
		Tag:         "udp4_addr",
		Translation: "{0} harus berupa alamat IPv4 UDP yang valid",
		Override:    false,
	},
	{
		Tag:         "udp6_addr",
		Translation: "{0} harus berupa alamat IPv6 UDP yang valid",
		Override:    false,
	},
	{
		Tag:         "ip_addr",
		Translation: "{0} harus berupa alamat IP yang dapat dipecahkan",
		Override:    false,
	},
	{
		Tag:         "ip4_addr",
		Translation: "{0} harus berupa alamat IPv4 yang dapat diatasi",
		Override:    false,
	},
	{
		Tag:         "ip6_addr",
		Translation: "{0} harus berupa alamat IPv6 yang dapat diatasi",
		Override:    false,
	},
	{
		Tag:         "unix_addr",
		Translation: "{0} harus berupa alamat UNIX yang dapat diatasi",
		Override:    false,
	},
	{
		Tag:         "mac",
		Translation: "{0} harus berisi alamat MAC yang valid",
		Override:    false,
	},
	{
		Tag:         "iscolor",
		Translation: "{0} harus berupa warna yang valid",
		Override:    false,
	},
	{
		Tag:         "oneof",
		Translation: "{0} harus berupa salah satu dari [{1}]",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
			s, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}
			return s
		},
	},
	{
		Tag:         "datetime",
		Translation: "{0} tidak sesuai format {1}",
		Override:    false,
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				log.Printf("warning: error translating FieldError: %#v", fe)
				return fe.(error).Error()
			}

			return t
		},
	},
	{
		Tag:         "boolean",
		Translation: "{0} tidak sesuai format boolean",
		Override:    false,
	},
	{
		Tag:         "image",
		Translation: "{0} harus berupa gambar yang valid",
		Override:    false,
	},
}
