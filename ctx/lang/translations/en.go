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

// reference: https://github.com/go-playground/validator/blob/master/translations/en/en.go
var English []translation = []translation{
	{
		Tag:         "required",
		Translation: "{0} is a required field",
		Override:    false,
	},
	{
		Tag:         "required_if",
		Translation: "{0} is a required field",
		Override:    false,
	},
	{
		Tag:         "required_unless",
		Translation: "{0} is a required field",
		Override:    false,
	},
	{
		Tag:         "required_with",
		Translation: "{0} is a required field",
		Override:    false,
	},
	{
		Tag:         "required_with_all",
		Translation: "{0} is a required field",
		Override:    false,
	},
	{
		Tag:         "required_without",
		Translation: "{0} is a required field",
		Override:    false,
	},
	{
		Tag:         "required_without_all",
		Translation: "{0} is a required field",
		Override:    false,
	},
	{
		Tag:         "excluded_if",
		Translation: "{0} is an excluded field",
		Override:    false,
	},
	{
		Tag:         "excluded_unless",
		Translation: "{0} is an excluded field",
		Override:    false,
	},
	{
		Tag:         "excluded_with",
		Translation: "{0} is an excluded field",
		Override:    false,
	},
	{
		Tag:         "excluded_with_all",
		Translation: "{0} is an excluded field",
		Override:    false,
	},
	{
		Tag:         "excluded_without",
		Translation: "{0} is an excluded field",
		Override:    false,
	},
	{
		Tag:         "excluded_without_all",
		Translation: "{0} is an excluded field",
		Override:    false,
	},
	{
		Tag:         "isdefault",
		Translation: "{0} must be default value",
		Override:    false,
	},
	{
		Tag: "len",
		CustomRegisFunc: func(ut ut.Translator) (err error) {
			if err = ut.Add("len-string", "{0} must be {1} in length", false); err != nil {
				return
			}

			if err = ut.AddCardinal("len-string-character", "{0} character", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("len-string-character", "{0} characters", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("len-number", "{0} must be equal to {1}", false); err != nil {
				return
			}

			if err = ut.Add("len-items", "{0} must contain {1}", false); err != nil {
				return
			}
			if err = ut.AddCardinal("len-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("len-items-item", "{0} items", locales.PluralRuleOther, false); err != nil {
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
			if err = ut.Add("min-string", "{0} must be at least {1} in length", false); err != nil {
				return
			}

			if err = ut.AddCardinal("min-string-character", "{0} character", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("min-string-character", "{0} characters", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("min-number", "{0} must be {1} or greater", false); err != nil {
				return
			}

			if err = ut.Add("min-items", "{0} must contain at least {1}", false); err != nil {
				return
			}
			if err = ut.AddCardinal("min-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("min-items-item", "{0} items", locales.PluralRuleOther, false); err != nil {
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
			if err = ut.Add("max-string", "{0} must be a maximum of {1} in length", false); err != nil {
				return
			}

			if err = ut.AddCardinal("max-string-character", "{0} character", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("max-string-character", "{0} characters", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("max-number", "{0} must be {1} or less", false); err != nil {
				return
			}

			if err = ut.Add("max-items", "{0} must contain at maximum {1}", false); err != nil {
				return
			}
			if err = ut.AddCardinal("max-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("max-items-item", "{0} items", locales.PluralRuleOther, false); err != nil {
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
		Translation: "{0} is not equal to {1}",
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
		Translation: "{0} should not be equal to {1}",
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
			if err = ut.Add("lt-string", "{0} must be less than {1} in length", false); err != nil {
				return
			}

			if err = ut.AddCardinal("lt-string-character", "{0} character", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("lt-string-character", "{0} characters", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lt-number", "{0} must be less than {1}", false); err != nil {
				return
			}

			if err = ut.Add("lt-items", "{0} must contain less than {1}", false); err != nil {
				return
			}

			if err = ut.AddCardinal("lt-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("lt-items-item", "{0} items", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lt-datetime", "{0} must be less than the current Date & Time", false); err != nil {
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
			if err = ut.Add("lte-string", "{0} must be at maximum {1} in length", false); err != nil {
				return
			}

			if err = ut.AddCardinal("lte-string-character", "{0} character", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("lte-string-character", "{0} characters", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lte-number", "{0} must be {1} or less", false); err != nil {
				return
			}

			if err = ut.Add("lte-items", "{0} must contain at maximum {1}", false); err != nil {
				return
			}

			if err = ut.AddCardinal("lte-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("lte-items-item", "{0} items", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("lte-datetime", "{0} must be less than or equal to the current Date & Time", false); err != nil {
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
			if err = ut.Add("gt-string", "{0} must be greater than {1} in length", false); err != nil {
				return
			}

			if err = ut.AddCardinal("gt-string-character", "{0} character", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("gt-string-character", "{0} characters", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gt-number", "{0} must be greater than {1}", false); err != nil {
				return
			}

			if err = ut.Add("gt-items", "{0} must contain more than {1}", false); err != nil {
				return
			}

			if err = ut.AddCardinal("gt-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("gt-items-item", "{0} items", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gt-datetime", "{0} must be greater than the current Date & Time", false); err != nil {
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
			if err = ut.Add("gte-string", "{0} must be at least {1} in length", false); err != nil {
				return
			}

			if err = ut.AddCardinal("gte-string-character", "{0} character", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("gte-string-character", "{0} characters", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gte-number", "{0} must be {1} or greater", false); err != nil {
				return
			}

			if err = ut.Add("gte-items", "{0} must contain at least {1}", false); err != nil {
				return
			}

			if err = ut.AddCardinal("gte-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
				return
			}

			if err = ut.AddCardinal("gte-items-item", "{0} items", locales.PluralRuleOther, false); err != nil {
				return
			}

			if err = ut.Add("gte-datetime", "{0} must be greater than or equal to the current Date & Time", false); err != nil {
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
		Translation: "{0} must be equal to {1}",
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
		Translation: "{0} must be equal to {1}",
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
		Translation: "{0} cannot be equal to {1}",
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
		Translation: "{0} must be greater than {1}",
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
		Translation: "{0} must be greater than or equal to {1}",
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
		Translation: "{0} must be less than {1}",
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
		Translation: "{0} must be less than or equal to {1}",
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
		Translation: "{0} cannot be equal to {1}",
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
		Translation: "{0} must be greater than {1}",
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
		Translation: "{0} must be greater than or equal to {1}",
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
		Translation: "{0} must be less than {1}",
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
		Translation: "{0} must be less than or equal to {1}",
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
		Translation: "{0} can only contain alphabetic characters",
		Override:    false,
	},
	{
		Tag:         "alphanum",
		Translation: "{0} can only contain alphanumeric characters",
		Override:    false,
	},
	{
		Tag:         "numeric",
		Translation: "{0} must be a valid numeric value",
		Override:    false,
	},
	{
		Tag:         "number",
		Translation: "{0} must be a valid number",
		Override:    false,
	},
	{
		Tag:         "hexadecimal",
		Translation: "{0} must be a valid hexadecimal",
		Override:    false,
	},
	{
		Tag:         "hexcolor",
		Translation: "{0} must be a valid HEX color",
		Override:    false,
	},
	{
		Tag:         "rgb",
		Translation: "{0} must be a valid RGB color",
		Override:    false,
	},
	{
		Tag:         "rgba",
		Translation: "{0} must be a valid RGBA color",
		Override:    false,
	},
	{
		Tag:         "hsl",
		Translation: "{0} must be a valid HSL color",
		Override:    false,
	},
	{
		Tag:         "hsla",
		Translation: "{0} must be a valid HSLA color",
		Override:    false,
	},
	{
		Tag:         "e164",
		Translation: "{0} must be a valid E.164 formatted phone number",
		Override:    false,
	},
	{
		Tag:         "email",
		Translation: "{0} must be a valid email address",
		Override:    false,
	},
	{
		Tag:         "url",
		Translation: "{0} must be a valid URL",
		Override:    false,
	},
	{
		Tag:         "uri",
		Translation: "{0} must be a valid URI",
		Override:    false,
	},
	{
		Tag:         "base64",
		Translation: "{0} must be a valid Base64 string",
		Override:    false,
	},
	{
		Tag:         "contains",
		Translation: "{0} must contain the text '{1}'",
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
		Translation: "{0} must contain at least one of the following characters '{1}'",
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
		Translation: "{0} cannot contain the text '{1}'",
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
		Translation: "{0} cannot contain any of the following characters '{1}'",
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
		Translation: "{0} cannot contain the following '{1}'",
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
		Translation: "{0} must be a valid ISBN number",
		Override:    false,
	},
	{
		Tag:         "isbn10",
		Translation: "{0} must be a valid ISBN-10 number",
		Override:    false,
	},
	{
		Tag:         "isbn13",
		Translation: "{0} must be a valid ISBN-13 number",
		Override:    false,
	},
	{
		Tag:         "issn",
		Translation: "{0} must be a valid ISSN number",
		Override:    false,
	},
	{
		Tag:         "uuid",
		Translation: "{0} must be a valid UUID",
		Override:    false,
	},
	{
		Tag:         "uuid3",
		Translation: "{0} must be a valid version 3 UUID",
		Override:    false,
	},
	{
		Tag:         "uuid4",
		Translation: "{0} must be a valid version 4 UUID",
		Override:    false,
	},
	{
		Tag:         "uuid5",
		Translation: "{0} must be a valid version 5 UUID",
		Override:    false,
	},
	{
		Tag:         "ulid",
		Translation: "{0} must be a valid ULID",
		Override:    false,
	},
	{
		Tag:         "ascii",
		Translation: "{0} must contain only ascii characters",
		Override:    false,
	},
	{
		Tag:         "printascii",
		Translation: "{0} must contain only printable ascii characters",
		Override:    false,
	},
	{
		Tag:         "multibyte",
		Translation: "{0} must contain multibyte characters",
		Override:    false,
	},
	{
		Tag:         "datauri",
		Translation: "{0} must contain a valid Data URI",
		Override:    false,
	},
	{
		Tag:         "latitude",
		Translation: "{0} must contain valid latitude coordinates",
		Override:    false,
	},
	{
		Tag:         "longitude",
		Translation: "{0} must contain a valid longitude coordinates",
		Override:    false,
	},
	{
		Tag:         "ssn",
		Translation: "{0} must be a valid SSN number",
		Override:    false,
	},
	{
		Tag:         "ipv4",
		Translation: "{0} must be a valid IPv4 address",
		Override:    false,
	},
	{
		Tag:         "ipv6",
		Translation: "{0} must be a valid IPv6 address",
		Override:    false,
	},
	{
		Tag:         "ip",
		Translation: "{0} must be a valid IP address",
		Override:    false,
	},
	{
		Tag:         "cidr",
		Translation: "{0} must contain a valid CIDR notation",
		Override:    false,
	},
	{
		Tag:         "cidrv4",
		Translation: "{0} must contain a valid CIDR notation for an IPv4 address",
		Override:    false,
	},
	{
		Tag:         "cidrv6",
		Translation: "{0} must contain a valid CIDR notation for an IPv6 address",
		Override:    false,
	},
	{
		Tag:         "tcp_addr",
		Translation: "{0} must be a valid TCP address",
		Override:    false,
	},
	{
		Tag:         "tcp4_addr",
		Translation: "{0} must be a valid IPv4 TCP address",
		Override:    false,
	},
	{
		Tag:         "tcp6_addr",
		Translation: "{0} must be a valid IPv6 TCP address",
		Override:    false,
	},
	{
		Tag:         "udp_addr",
		Translation: "{0} must be a valid UDP address",
		Override:    false,
	},
	{
		Tag:         "udp4_addr",
		Translation: "{0} must be a valid IPv4 UDP address",
		Override:    false,
	},
	{
		Tag:         "udp6_addr",
		Translation: "{0} must be a valid IPv6 UDP address",
		Override:    false,
	},
	{
		Tag:         "ip_addr",
		Translation: "{0} must be a resolvable IP address",
		Override:    false,
	},
	{
		Tag:         "ip4_addr",
		Translation: "{0} must be a resolvable IPv4 address",
		Override:    false,
	},
	{
		Tag:         "ip6_addr",
		Translation: "{0} must be a resolvable IPv6 address",
		Override:    false,
	},
	{
		Tag:         "unix_addr",
		Translation: "{0} must be a resolvable UNIX address",
		Override:    false,
	},
	{
		Tag:         "mac",
		Translation: "{0} must contain a valid MAC address",
		Override:    false,
	},
	{
		Tag:         "fqdn",
		Translation: "{0} must be a valid FQDN",
		Override:    false,
	},
	{
		Tag:         "unique",
		Translation: "{0} must contain unique values",
		Override:    false,
	},
	{
		Tag:         "iscolor",
		Translation: "{0} must be a valid color",
		Override:    false,
	},
	{
		Tag:         "cron",
		Translation: "{0} must be a valid cron expression",
		Override:    false,
	},
	{
		Tag:         "oneof",
		Translation: "{0} must be one of [{1}]",
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
		Tag:         "json",
		Translation: "{0} must be a valid json string",
		Override:    false,
	},
	{
		Tag:         "jwt",
		Translation: "{0} must be a valid jwt string",
		Override:    false,
	},
	{
		Tag:         "lowercase",
		Translation: "{0} must be a lowercase string",
		Override:    false,
	},
	{
		Tag:         "uppercase",
		Translation: "{0} must be an uppercase string",
		Override:    false,
	},
	{
		Tag:         "datetime",
		Translation: "{0} does not match the {1} format",
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
		Translation: "{0} does not match postcode format of {1} country",
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
		Translation: "{0} does not match postcode format of country in {1} field",
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
		Translation: "{0} must be a valid boolean value",
		Override:    false,
	},
	{
		Tag:         "image",
		Translation: "{0} must be a valid image",
		Override:    false,
	},
	{
		Tag:         "cve",
		Translation: "{0} must be a valid cve identifier",
		Override:    false,
	},
}
