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
	"strings"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	transLang "github.com/dedyf5/resik/ctx/lang/translations"
	"github.com/dedyf5/resik/pkg/status"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	idTranslation "github.com/go-playground/validator/v10/translations/id"
	"golang.org/x/text/language"
)

//go:generate mockgen -source validator.go -package mock -destination ./mock/validator.go
type IValidate interface {
	Struct(payloadStruct interface{}, lang *langCtx.Lang) *status.Status
	ErrorFormatter(err error, lang *langCtx.Lang) *status.Status
}

type Validate struct {
	instance            *validator.Validate
	langDefault         language.Tag
	translatorDefault   ut.Translator
	universalTranslator *ut.UniversalTranslator
}

func New(langDefault language.Tag) *Validate {
	validate := validator.New()
	uni := ut.New(transLang.LanguageToTranslator(langDefault), transLang.Translators()...)
	trans, found := uni.GetTranslator(langDefault.String())
	if !found {
		log.Panic("translator not found")
	}

	registerDefaultTranslations(langDefault, validate, trans)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	for lang, translations := range transLang.Map {
		if lang == langDefault {
			continue
		}
		engine, _ := uni.FindTranslator(lang.String())
		for _, translation := range translations {
			var err error = nil
			if translation.CustomTransFunc != nil && translation.CustomRegisFunc != nil {
				err = validate.RegisterTranslation(translation.Tag, engine, translation.CustomRegisFunc, translation.CustomTransFunc)
			} else if translation.CustomTransFunc != nil && translation.CustomRegisFunc == nil {
				err = validate.RegisterTranslation(translation.Tag, engine, registrationFunc(translation.Tag, translation.Translation, translation.Override), translation.CustomTransFunc)
			} else if translation.CustomTransFunc == nil && translation.CustomRegisFunc != nil {
				err = validate.RegisterTranslation(translation.Tag, engine, translation.CustomRegisFunc, translateFunc)
			} else {
				err = validate.RegisterTranslation(translation.Tag, engine, registrationFunc(translation.Tag, translation.Translation, translation.Override), translateFunc)
			}
			if err != nil {
				log.Panicf("register translation failed (lang: %s) %s", lang.String(), err.Error())
			}
		}
	}

	return &Validate{
		instance:            validate,
		langDefault:         langDefault,
		translatorDefault:   trans,
		universalTranslator: uni,
	}
}

func (v *Validate) Struct(payloadStruct interface{}, lang *langCtx.Lang) *status.Status {
	if payloadStruct == nil {
		return &status.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errors.New("payloadStruct can't be nil"),
		}
	}

	err := v.instance.Struct(payloadStruct)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return &status.Status{
				Code: http.StatusInternalServerError,
			}
		}
		return v.ErrorFormatter(err, lang)
	}

	return nil
}

func (v *Validate) ErrorFormatter(err error, lang *langCtx.Lang) *status.Status {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return &status.Status{
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

		switch e.Tag() {
		case "datetime":
			{
				value = regexDate.ReplaceAllString(value, "yyyy-MM-dd")
				value = regexTime.ReplaceAllString(value, "HH:mm:ss")
			}
		}

		errMap[field] = value
		if k == 0 {
			first = value
		}
	}

	return &status.Status{
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

func getLanguage(langDef language.Tag, lang *langCtx.Lang) language.Tag {
	if lang == nil {
		return langDef
	}
	return lang.LanguageReqOrDefault()
}

func registerDefaultTranslations(lang language.Tag, validate *validator.Validate, trans ut.Translator) {
	switch lang.String() {
	case "id":
		_ = idTranslation.RegisterDefaultTranslations(validate, trans)
	default:
		enTranslation.RegisterDefaultTranslations(validate, trans)
	}
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		log.Printf("warning: error translating FieldError: %#v", fe)
		return fe.(error).Error()
	}

	return t
}
