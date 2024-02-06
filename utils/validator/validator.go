// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package validator

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	transLang "github.com/dedyf5/resik/ctx/lang/translations"
	resPkg "github.com/dedyf5/resik/pkg/response"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
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

	uni := ut.New(transLang.LanguageToTranslator(langDefault), transLang.Translators()...)
	trans, found := uni.GetTranslator(langDefault.String())
	if !found {
		log.Panic("translator not found")
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	for lang, translations := range transLang.Map {
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

func (v *Validate) Struct(payloadStruct interface{}, lang *langCtx.Lang) *resPkg.Status {
	if payloadStruct == nil {
		return &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errors.New("payloadStruct can't be nil"),
		}
	}

	err := v.instance.Struct(payloadStruct)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return &resPkg.Status{
				Code: http.StatusInternalServerError,
			}
		}
		return v.ErrorFormatter(err, lang)
	}

	return nil
}

func (v *Validate) ErrorFormatter(err error, lang *langCtx.Lang) *resPkg.Status {
	errs, ok := err.(validator.ValidationErrors)
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
	s = strings.ReplaceAll(s, `"`, ``)
	base := strings.Split(s, " ")
	vals := make([]string, 0, cap(base)*2)
	for _, v := range base {
		vals = append(vals, v)
		vals = append(vals, fmt.Sprintf("-%s", v))
	}

	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
	for i := 0; i < len(vals); i++ {
		if vals[i] == v {
			return true
		}
	}
	return false
}

func getLanguage(langDef language.Tag, lang *langCtx.Lang) language.Tag {
	if lang == nil {
		return langDef
	}
	return lang.LanguageReqOrDefault()
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
