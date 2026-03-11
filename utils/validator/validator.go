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
	enTrans "github.com/go-playground/validator/v10/translations/en"
	idTrans "github.com/go-playground/validator/v10/translations/id"
	jaTrans "github.com/go-playground/validator/v10/translations/ja"
	"golang.org/x/text/language"
)

//go:generate mockgen -source validator.go -package mock -destination ./mock/validator.go
type IValidate interface {
	Struct(payloadStruct interface{}, lang *langCtx.Lang) *resPkg.Status
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

	if err := validate.RegisterValidation("oneof_order", isOneOfOrder); err != nil {
		log.Panic("error register validation tag oneof_order")
	}

	uni := ut.New(LanguageToTranslator(langDefault), Translators()...)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonTag := fld.Tag.Get("json")

		if jsonTag != "" && jsonTag != "-" {
			for name := range strings.SplitSeq(jsonTag, ",") {
				return name
			}
		}

		protoTag := fld.Tag.Get("protobuf")
		if protoTag != "" {
			for p := range strings.SplitSeq(protoTag, ",") {
				if after, ok := strings.CutPrefix(p, "name="); ok {
					return after
				}
			}
		}

		return fld.Name
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
		if err := registerOneOfOrder(v, t, "{0} harus berupa salah satu dari [{1}]"); err != nil {
			log.Printf("[validator] failed to register oneof_order for ID: %v", err)
		}
	}

	// English
	if t, found := uni.GetTranslator(language.English.String()); found {
		if err := enTrans.RegisterDefaultTranslations(v, t); err != nil {
			log.Printf("[validator] failed to register EN translation: %v", err)
		}
		if err := registerOneOfOrder(v, t, "{0} must be one of [{1}]"); err != nil {
			log.Printf("[validator] failed to register oneof_order for EN: %v", err)
		}
	}

	// Japanese
	if t, found := uni.GetTranslator(language.Japanese.String()); found {
		if err := jaTrans.RegisterDefaultTranslations(v, t); err != nil {
			log.Printf("[validator] failed to register JA translation: %v", err)
		}
		if err := registerOneOfOrder(v, t, "{0}は[{1}]のうちのいずれかでなければなりません"); err != nil {
			log.Printf("[validator] failed to register oneof_order for JA: %v", err)
		}
	}
}

func registerOneOfOrder(v *validator.Validate, trans ut.Translator, msg string) error {
	return v.RegisterTranslation("oneof_order", trans, func(ut ut.Translator) error {
		return ut.Add("oneof_order", msg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		s := strings.ReplaceAll(fe.Param(), "'", "")
		s = strings.ReplaceAll(s, `\"`, ``)
		base := strings.Split(s, " ")
		vals := make([]string, 0, len(base)*2)
		for _, v := range base {
			vals = append(vals, v)
			vals = append(vals, "-"+v)
		}
		param := strings.Join(vals, " ")

		t, _ := ut.T("oneof_order", fe.Field(), param)
		return t
	})
}

func (v *Validate) Struct(payloadStruct interface{}, lang *langCtx.Lang) *resPkg.Status {
	if payloadStruct == nil {
		return &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errors.New("payloadStruct can't be nil"),
		}
	}

	err := v.instance.Struct(payloadStruct)
	if err != nil {
		if _, ok := errors.AsType[*validator.InvalidValidationError](err); ok {
			return &resPkg.Status{
				Code: http.StatusInternalServerError,
			}
		}
		return v.ErrorFormatter(err, lang)
	}

	return nil
}

func (v *Validate) ErrorFormatter(err error, lang *langCtx.Lang) *resPkg.Status {
	errs, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		return &resPkg.Status{
			Code: http.StatusBadRequest,
		}
	}

	errMap := map[string]string{}

	regexKey := regexp.MustCompile(`^[^.]*.`)
	regexDate := regexp.MustCompile("2006-01-02")
	regexTime := regexp.MustCompile("15:04:05")
	first := ""
	for k, e := range errs {
		// replace nested field name
		field := regexKey.ReplaceAllString(e.Namespace(), "")

		// replace field name in error message
		value := e.Translate(v.Translator(getLanguage(v.langDefault, lang)))
		value = regexp.MustCompile(e.Field()).ReplaceAllString(value, field)

		if e.Tag() == "datetime" {
			value = regexDate.ReplaceAllString(value, "yyyy-MM-dd")
			value = regexTime.ReplaceAllString(value, "HH:mm:ss")
		}

		errMap[field] = value
		if k == 0 {
			first = value
		}
	}

	return &resPkg.Status{
		Code:    http.StatusBadRequest,
		Message: first,
		Detail:  errMap,
	}
}

func (v *Validate) Translator(lang language.Tag) ut.Translator {
	t, found := v.universalTranslator.GetTranslator(lang.String())
	if !found {
		t = v.translatorDefault
	}
	return t
}

func isOneOfOrder(fl validator.FieldLevel) bool {
	s := strings.ReplaceAll(fl.Param(), "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	base := strings.Split(s, " ")
	vals := make([]string, 0, len(base)*2)
	for _, v := range base {
		vals = append(vals, v)
		vals = append(vals, "-"+v)
	}

	fields := strings.SplitSeq(fl.Field().String(), ",")
	for field := range fields {
		if !slices.Contains(vals, field) {
			return false
		}
	}

	return true
}

func getLanguage(langDef language.Tag, lang *langCtx.Lang) language.Tag {
	if lang == nil {
		return langDef
	}
	return lang.LanguageReqOrDefault()
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
