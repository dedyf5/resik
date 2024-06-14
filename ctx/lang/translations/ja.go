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

// reference: https://github.com/go-playground/validator/blob/master/translations/ja/ja.go
var Japanese []translation = []translation{
	{
		Tag:         "required",
		Translation: "{0}は必須フィールドです",
		Override:    false,
	},
	{
		Tag:         "required_if",
		Translation: "{0}は必須フィールドです",
		Override:    false,
	},
	{
		Tag: "len",
		CustomRegisFunc: func(ut ut.Translator) (err error) {
			if err = ut.Add("len-string", "{0}の長さは{1}でなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("len-string-character", "{0}文字", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("len-string-character", "{0}文字", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("len-number", "{0}は{1}と等しくなければなりません", false); err != nil {
				return
			}

			if err = ut.Add("len-items", "{0}は{1}を含まなければなりません", false); err != nil {
				return
			}
			// if err = ut.AddCardinal("len-items-item", "{0}つの項目", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("len-items-item", "{0}つの項目", locales.PluralRuleOther, false); err != nil {
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
			if err = ut.Add("min-string", "{0}の長さは少なくとも{1}はなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("min-string-character", "{0}文字", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("min-string-character", "{0}文字", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("min-number", "{0}は{1}以上でなければなりません", false); err != nil {
				return
			}

			if err = ut.Add("min-items", "{0}は少なくとも{1}を含まなければなりません", false); err != nil {
				return
			}
			// if err = ut.AddCardinal("min-items-item", "{0}つの項目", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("min-items-item", "{0}つの項目", locales.PluralRuleOther, false); err != nil {
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
			if err = ut.Add("max-string", "{0}の長さは最大でも{1}でなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("max-string-character", "{0}文字", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("max-string-character", "{0}文字", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("max-number", "{0}は{1}以下でなければなりません", false); err != nil {
				return
			}

			if err = ut.Add("max-items", "{0}は最大でも{1}でなければなりません", false); err != nil {
				return
			}
			// if err = ut.AddCardinal("max-items-item", "{0}つの項目", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("max-items-item", "{0}つの項目", locales.PluralRuleOther, false); err != nil {
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
		Translation: "{0}は{1}と等しくありません",
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
		Tag: "ne",
		CustomRegisFunc: func(ut ut.Translator) (err error) {
			if err = ut.Add("ne-items", "{0}の項目の数は{1}と異ならなければなりません", false); err != nil {
				fmt.Printf("ne customRegisFunc #1 error because of %v\n", err)
				return
			}
			// if err = ut.AddCardinal("ne-items-item", "{0}個", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("ne-items-item", "{0}個", locales.PluralRuleOther, false); err != nil {
				return
			}
			if err = ut.Add("ne", "{0}は{1}と異ならなければなりません", false); err != nil {
				fmt.Printf("ne customRegisFunc #2 error because of %v\n", err)
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
			case reflect.Slice:
				var c string
				err = fn()
				if err != nil {
					goto END
				}

				c, err = ut.C("ne-items-item", f64, digits, ut.FmtNumber(f64, digits))
				if err != nil {
					goto END
				}
				t, err = ut.T("ne-items", fe.Field(), c)
			default:
				t, err = ut.T("ne", fe.Field(), fe.Param())
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
		Tag: "lt",
		CustomRegisFunc: func(ut ut.Translator) (err error) {
			if err = ut.Add("lt-string", "{0}の長さは{1}よりも少なくなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("lt-string-character", "{0}文字", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("lt-string-character", "{0}文字", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lt-number", "{0}は{1}よりも小さくなければなりません", false); err != nil {
				return
			}

			if err = ut.Add("lt-items", "{0}は{1}よりも少ない項目でなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("lt-items-item", "{0}つの項目", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("lt-items-item", "{0}つの項目", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lt-datetime", "{0}は現時刻よりも前でなければなりません", false); err != nil {
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
			if err = ut.Add("lte-string", "{0}の長さは最大でも{1}でなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("lte-string-character", "{0}文字", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("lte-string-character", "{0}文字", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lte-number", "{0}は{1}以下でなければなりません", false); err != nil {
				return
			}

			if err = ut.Add("lte-items", "{0}は最大でも{1}でなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("lte-items-item", "{0}つの項目", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("lte-items-item", "{0}つの項目", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lte-datetime", "{0}は現時刻以前でなければなりません", false); err != nil {
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
			if err = ut.Add("gt-string", "{0}の長さは{1}よりも多くなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("gt-string-character", "{0}文字", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("gt-string-character", "{0}文字", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gt-number", "{0}は{1}よりも大きくなければなりません", false); err != nil {
				return
			}

			if err = ut.Add("gt-items", "{0}は{1}よりも多い項目を含まなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("gt-items-item", "{0}つの項目", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("gt-items-item", "{0}つの項目", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gt-datetime", "{0}は現時刻よりも後でなければなりません", false); err != nil {
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
			if err = ut.Add("gte-string", "{0}の長さは少なくとも{1}以上はなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("gte-string-character", "{0}文字", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("gte-string-character", "{0}文字", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gte-number", "{0}は{1}以上でなければなりません", false); err != nil {
				return
			}

			if err = ut.Add("gte-items", "{0}は少なくとも{1}を含まなければなりません", false); err != nil {
				return
			}

			// if err = ut.AddCardinal("gte-items-item", "{0}つの項目", locales.PluralRuleOne, false); err != nil {
			// 	return
			// }

			if err = ut.AddCardinal("gte-items-item", "{0}つの項目", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gte-datetime", "{0}は現時刻以降でなければなりません", false); err != nil {
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
		Translation: "{0}は{1}と等しくなければなりません",
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
		Translation: "{0}は{1}と等しくなければなりません",
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
		Translation: "{0}は{1}とは異ならなければなりません",
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
		Translation: "{0}は{1}よりも大きくなければなりません",
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
		Translation: "{0}は{1}以上でなければなりません",
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
		Translation: "{0}は{1}よりも小さくなければなりません",
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
		Translation: "{0}は{1}以下でなければなりません",
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
		Translation: "{0}は{1}とは異ならなければなりません",
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
		Translation: "{0}は{1}よりも大きくなければなりません",
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
		Translation: "{0}は{1}以上でなければなりません",
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
		Translation: "{0}は{1}よりも小さくなければなりません",
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
		Translation: "{0}は{1}以下でなければなりません",
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
		Translation: "{0}はアルファベットのみを含むことができます",
		Override:    false,
	},
	{
		Tag:         "alphanum",
		Translation: "{0}はアルファベットと数字のみを含むことができます",
		Override:    false,
	},
	{
		Tag:         "numeric",
		Translation: "{0}は正しい数字でなければなりません",
		Override:    false,
	},
	{
		Tag:         "number",
		Translation: "{0}は正しい数でなければなりません",
		Override:    false,
	},
	{
		Tag:         "hexadecimal",
		Translation: "{0}は正しい16進表記でなければなりません",
		Override:    false,
	},
	{
		Tag:         "hexcolor",
		Translation: "{0}は正しいHEXカラーコードでなければなりません",
		Override:    false,
	},
	{
		Tag:         "rgb",
		Translation: "{0}は正しいRGBカラーコードでなければなりません",
		Override:    false,
	},
	{
		Tag:         "rgba",
		Translation: "{0}は正しいRGBAカラーコードでなければなりません",
		Override:    false,
	},
	{
		Tag:         "hsl",
		Translation: "{0}は正しいHSLカラーコードでなければなりません",
		Override:    false,
	},
	{
		Tag:         "hsla",
		Translation: "{0}は正しいHSLAカラーコードでなければなりません",
		Override:    false,
	},
	{
		Tag:         "email",
		Translation: "{0}は正しいメールアドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "url",
		Translation: "{0}は正しいURLでなければなりません",
		Override:    false,
	},
	{
		Tag:         "uri",
		Translation: "{0}は正しいURIでなければなりません",
		Override:    false,
	},
	{
		Tag:         "base64",
		Translation: "{0}は正しいBase64文字列でなければなりません",
		Override:    false,
	},
	{
		Tag:         "contains",
		Translation: "{0}は'{1}'を含まなければなりません",
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
		Translation: "{0}は'{1}'の少なくとも1つを含まなければなりません",
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
		Translation: "{0}には'{1}'というテキストを含むことはできません",
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
		Translation: "{0}には'{1}'のどれも含めることはできません",
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
		Translation: "{0}には'{1}'を含めることはできません",
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
		Translation: "{0}は正しいISBN番号でなければなりません",
		Override:    false,
	},
	{
		Tag:         "isbn10",
		Translation: "{0}は正しいISBN-10番号でなければなりません",
		Override:    false,
	},
	{
		Tag:         "isbn13",
		Translation: "{0}は正しいISBN-13番号でなければなりません",
		Override:    false,
	},
	{
		Tag:         "issn",
		Translation: "{0}は正しいISSN番号でなければなりません",
		Override:    false,
	},
	{
		Tag:         "uuid",
		Translation: "{0}は正しいUUIDでなければなりません",
		Override:    false,
	},
	{
		Tag:         "uuid3",
		Translation: "{0}はバージョンが3の正しいUUIDでなければなりません",
		Override:    false,
	},
	{
		Tag:         "uuid4",
		Translation: "{0}はバージョンが4の正しいUUIDでなければなりません",
		Override:    false,
	},
	{
		Tag:         "uuid5",
		Translation: "{0}はバージョンが5の正しいUUIDでなければなりません",
		Override:    false,
	},
	{
		Tag:         "ulid",
		Translation: "{0}は正しいULIDでなければなりません",
		Override:    false,
	},
	{
		Tag:         "ascii",
		Translation: "{0}はASCII文字のみを含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "printascii",
		Translation: "{0}は印刷可能なASCII文字のみを含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "multibyte",
		Translation: "{0}はマルチバイト文字を含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "datauri",
		Translation: "{0}は正しいデータURIを含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "latitude",
		Translation: "{0}は正しい緯度の座標を含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "longitude",
		Translation: "{0}は正しい経度の座標を含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "ssn",
		Translation: "{0}は正しい社会保障番号でなければなりません",
		Override:    false,
	},
	{
		Tag:         "ipv4",
		Translation: "{0}は正しいIPv4アドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "ipv6",
		Translation: "{0}は正しいIPv6アドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "ip",
		Translation: "{0}は正しいIPアドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "cidr",
		Translation: "{0}は正しいCIDR表記を含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "cidrv4",
		Translation: "{0}はIPv4アドレスの正しいCIDR表記を含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "cidrv6",
		Translation: "{0}はIPv6アドレスの正しいCIDR表記を含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "tcp_addr",
		Translation: "{0}は正しいTCPアドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "tcp4_addr",
		Translation: "{0}は正しいIPv4のTCPアドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "tcp6_addr",
		Translation: "{0}は正しいIPv6のTCPアドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "udp_addr",
		Translation: "{0}は正しいUDPアドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "udp4_addr",
		Translation: "{0}は正しいIPv4のUDPアドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "udp6_addr",
		Translation: "{0}は正しいIPv6のUDPアドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "ip_addr",
		Translation: "{0}は解決可能なIPアドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "ip4_addr",
		Translation: "{0}は解決可能なIPv4アドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "ip6_addr",
		Translation: "{0}は解決可能なIPv6アドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "unix_addr",
		Translation: "{0}は解決可能なUNIXアドレスでなければなりません",
		Override:    false,
	},
	{
		Tag:         "mac",
		Translation: "{0}は正しいMACアドレスを含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "unique",
		Translation: "{0}は一意な値のみを含まなければなりません",
		Override:    false,
	},
	{
		Tag:         "iscolor",
		Translation: "{0}は正しい色でなければなりません",
		Override:    false,
	},
	{
		Tag:         "oneof",
		Translation: "{0}は[{1}]のうちのいずれかでなければなりません",
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
		Tag:         "image",
		Translation: "{0} は有効な画像でなければなりません",
		Override:    false,
	},
	{
		Tag:         "json",
		Translation: "{0}は正しいJSON文字列でなければなりません",
		Override:    false,
	},
	{
		Tag:         "jwt",
		Translation: "{0}は正しいJWT文字列でなければなりません",
		Override:    false,
	},
	{
		Tag:         "lowercase",
		Translation: "{0}は小文字でなければなりません",
		Override:    false,
	},
	{
		Tag:         "uppercase",
		Translation: "{0}は大文字でなければなりません",
		Override:    false,
	},
	{
		Tag:         "datetime",
		Translation: "{0}は{1}の書式と一致しません",
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
		Tag:         "postcode_iso3166_alpha2",
		Translation: "{0}は国名コード{1}の郵便番号形式と一致しません",
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
		Tag:         "postcode_iso3166_alpha2_field",
		Translation: "{0}は{1}フィールドで指定された国名コードの郵便番号形式と一致しません",
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
		Translation: "{0}は正しいブール値でなければなりません",
		Override:    false,
	},
}
