// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package validator

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"slices"
	"strings"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	idTrans "github.com/go-playground/validator/v10/translations/id"
	jaTrans "github.com/go-playground/validator/v10/translations/ja"
	"golang.org/x/text/language"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

//go:generate mockgen -source validator.go -package mock -destination ./mock/validator.go
type IValidate interface {
	Struct(payloadStruct any, lang *langCtx.Lang) *resPkg.Status
	ErrorFormatter(err error, lang *langCtx.Lang) *resPkg.Status
}

type Validate struct {
	instance            *validator.Validate
	langDefault         language.Tag
	translatorDefault   ut.Translator
	universalTranslator *ut.UniversalTranslator
}

func New(langDefault language.Tag) *Validate {
	validate := validator.New()

	if err := validate.RegisterValidation("notblank", validators.NotBlank); err != nil {
		log.Panic("error register validation tag notblank")
	}

	if err := validate.RegisterValidation("oneof_order", isOneOfOrder); err != nil {
		log.Panic("error register validation tag oneof_order")
	}

	uni := ut.New(LanguageToTranslator(langDefault), Translators()...)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonTag := fld.Tag.Get("json")

		if jsonTag != "" && jsonTag != "-" {
			for name := range strings.SplitSeq(jsonTag, ",") {
				return fieldPlaceholder(name)
			}
		}

		protoTag := fld.Tag.Get("protobuf")
		if protoTag != "" {
			for p := range strings.SplitSeq(protoTag, ",") {
				if after, ok := strings.CutPrefix(p, "name="); ok {
					return fieldPlaceholder(after)
				}
			}
		}

		return fieldPlaceholder(fld.Name)
	})

	registerAllTranslations(validate, uni)

	trans, _ := uni.GetTranslator(langDefault.String())

	return &Validate{
		instance:            validate,
		langDefault:         langDefault,
		translatorDefault:   trans,
		universalTranslator: uni,
	}
}

func registerAllTranslations(v *validator.Validate, uni *ut.UniversalTranslator) {
	// Indonesian
	if t, found := uni.GetTranslator(language.Indonesian.String()); found {
		if err := idTrans.RegisterDefaultTranslations(v, t); err != nil {
			log.Printf("[validator] failed to register ID translation: %v", err)
		}
		if err := registerNotBlank(v, t, "{0} tidak boleh kosong"); err != nil {
			log.Printf("[validator] failed to register notblank for ID: %v", err)
		}
		if err := registerOneOfOrder(v, t, "{0} harus berupa salah satu dari [{1}]"); err != nil {
			log.Printf("[validator] failed to register oneof_order for ID: %v", err)
		}
	}

	// English
	if t, found := uni.GetTranslator(language.English.String()); found {
		if err := enTrans.RegisterDefaultTranslations(v, t); err != nil {
			log.Printf("[validator] failed to register EN translation: %v", err)
		}
		if err := registerNotBlank(v, t, "{0} must not be blank"); err != nil {
			log.Printf("[validator] failed to register notblank for EN: %v", err)
		}
		if err := registerOneOfOrder(v, t, "{0} must be one of [{1}]"); err != nil {
			log.Printf("[validator] failed to register oneof_order for EN: %v", err)
		}

		err := v.RegisterTranslation("timezone", t, func(ut ut.Translator) error {
			return ut.Add("timezone", "{0} must be a valid timezone", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("timezone", fe.Field())
			return t
		})
		if err != nil {
			log.Printf("[validator] failed to register timezone translation for EN: %v", err)
		}
	}

	// Japanese
	if t, found := uni.GetTranslator(language.Japanese.String()); found {
		if err := jaTrans.RegisterDefaultTranslations(v, t); err != nil {
			log.Printf("[validator] failed to register JA translation: %v", err)
		}
		if err := registerNotBlank(v, t, "{0}は空であってはなりません"); err != nil {
			log.Printf("[validator] failed to register notblank for JA: %v", err)
		}
		if err := registerOneOfOrder(v, t, "{0}は[{1}]のうちのいずれかでなければなりません"); err != nil {
			log.Printf("[validator] failed to register oneof_order for JA: %v", err)
		}

		err := v.RegisterTranslation("timezone", t, func(ut ut.Translator) error {
			return ut.Add("timezone", "{0}は有効なタイムゾーンである必要があります", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("timezone", fe.Field())
			return t
		})
		if err != nil {
			log.Printf("[validator] failed to register timezone translation for JA: %v", err)
		}
	}
}

func registerNotBlank(v *validator.Validate, trans ut.Translator, msg string) error {
	return v.RegisterTranslation("notblank", trans, func(ut ut.Translator) error {
		return ut.Add("notblank", msg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notblank", fe.Field())
		return t
	})
}

func registerOneOfOrder(v *validator.Validate, trans ut.Translator, msg string) error {
	return v.RegisterTranslation("oneof_order", trans, func(ut ut.Translator) error {
		return ut.Add("oneof_order", msg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		options := getOneOfOrderOptions(fe.Param())
		param := strings.Join(options, " ")

		t, _ := ut.T("oneof_order", fe.Field(), param)
		return t
	})
}

func (v *Validate) Struct(payloadStruct any, lang *langCtx.Lang) *resPkg.Status {
	if payloadStruct == nil {
		return resPkg.NewStatusError(
			http.StatusInternalServerError,
			errors.New("payloadStruct can't be nil"),
		)
	}

	err := v.instance.Struct(payloadStruct)
	if err != nil {
		if _, ok := errors.AsType[*validator.InvalidValidationError](err); ok {
			return resPkg.NewStatusCode(http.StatusInternalServerError)
		}
		return v.ErrorFormatter(err, lang)
	}

	return nil
}

func (v *Validate) ErrorFormatter(err error, lang *langCtx.Lang) *resPkg.Status {
	errs, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		return resPkg.NewStatusCode(http.StatusBadRequest)
	}

	volations := []*errdetails.BadRequest_FieldViolation{}

	regexKey := regexp.MustCompile(`^[^.]*.`)
	regexDate := regexp.MustCompile("2006-01-02")
	regexTime := regexp.MustCompile("T15:04:05")
	regexZone := regexp.MustCompile("Z07:00")

	first := ""
	currLang := getLangReqOrDefault(v.langDefault, lang)
	locale := currLang.LanguageReqOrDefault()
	translator := v.Translator(locale)

	for k, e := range errs {
		rawFieldPath := removeFieldPlaceholder(regexKey.ReplaceAllString(e.Namespace(), ""))
		translatedFieldName := currLang.GetValidationFieldNameWithQuote(rawFieldPath)
		value := e.Translate(translator)
		value = strings.ReplaceAll(value, e.Field(), translatedFieldName)

		if e.Tag() == "datetime" {
			value = regexDate.ReplaceAllString(value, "yyyy-MM-dd")
			value = regexTime.ReplaceAllString(value, "'T'HH:mm:ss")
			value = regexZone.ReplaceAllString(value, ".SSSZZ")
		}

		volations = append(volations, &errdetails.BadRequest_FieldViolation{
			Field:       rawFieldPath,
			Description: e.Error(),
			Reason:      errorReason(e),
			LocalizedMessage: &errdetails.LocalizedMessage{
				Locale:  currLang.LanguageReqOrDefault().String(),
				Message: value,
			},
		})

		if k == 0 {
			first = value
		}
	}

	badRequest := &errdetails.BadRequest{
		FieldViolations: volations,
	}

	return resPkg.NewStatusDetails(http.StatusBadRequest, first, badRequest)
}

func (v *Validate) Translator(lang language.Tag) ut.Translator {
	t, found := v.universalTranslator.GetTranslator(lang.String())
	if !found {
		t = v.translatorDefault
	}
	return t
}

func isOneOfOrder(fl validator.FieldLevel) bool {
	options := getOneOfOrderOptions(fl.Param())
	fields := strings.SplitSeq(fl.Field().String(), ",")
	for field := range fields {
		if !slices.Contains(options, field) {
			return false
		}
	}

	return true
}

func getOneOfOrderOptions(param string) []string {
	s := strings.ReplaceAll(param, "'", "")
	s = strings.ReplaceAll(s, `\"`, ``)
	base := strings.Split(s, " ")
	vals := make([]string, 0, len(base)*2)
	for _, v := range base {
		vals = append(vals, v)
		vals = append(vals, "-"+v)
	}
	return vals
}

func errorReason(err validator.FieldError) string {
	switch err.Tag() {
	case "required", "required_if", "required_unless", "required_with", "required_with_all", "required_without", "required_without_all":
		return "REQUIRED"

	case "oneof", "oneof_order":
		return "OUT_OF_RANGE"

	case "datetime", "email", "url", "uuid", "hostname", "ip", "latitude", "longitude":
		return "INVALID_FORMAT"

	case "min", "gt", "gte":
		if err.Kind() == reflect.String || err.Kind() == reflect.Slice || err.Kind() == reflect.Map {
			return "TOO_SHORT"
		}
		return "TOO_SMALL"

	case "max", "lt", "lte":
		if err.Kind() == reflect.String || err.Kind() == reflect.Slice || err.Kind() == reflect.Map {
			return "TOO_LONG"
		}
		return "TOO_LARGE"

	case "len":
		return "INVALID_LENGTH"

	case "numeric", "number":
		return "TYPE_MISMATCH"
	case "alpha", "alphanum", "ascii":
		return "INVALID_FORMAT"
	case "eq", "ne":
		return "INVALID"
	case "eqfield", "nefield":
		return "FIELD_CONFLICT"

	default:
		return "INVALID"
	}
}

func fieldPlaceholder(field string) string {
	return "{" + field + "}"
}

func removeFieldPlaceholder(field string) string {
	return strings.TrimPrefix(strings.TrimSuffix(field, "}"), "{")
}

func getLangReqOrDefault(languageDefault language.Tag, lang *langCtx.Lang) *langCtx.Lang {
	if lang == nil {
		return langCtx.NewLang(languageDefault, nil, "")
	}
	return lang
}

func Translators() (res []locales.Translator) {
	for _, v := range langCtx.Available {
		res = append(res, LanguageToTranslator(v))
	}
	return
}

func LanguageToTranslator(lang language.Tag) locales.Translator {
	switch lang {
	case language.Indonesian:
		return id.New()
	case language.Japanese:
		return ja.New()
	}
	return en.New()
}
