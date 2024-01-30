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

	"github.com/dedyf5/resik/ctx/status"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	"golang.org/x/text/language"
)

//go:generate mockgen -source validator.go -package mock -destination ./mock/validator.go
type IValidate interface {
	Struct(payloadStruct interface{}) *status.Status
	ErrorFormatter(err error) *status.Status
}

type Validate struct {
	instance    *validator.Validate
	langDefault language.Tag
	translator  ut.Translator
}

func New(langDefault language.Tag) *Validate {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, found := uni.GetTranslator(langDefault.String())
	if !found {
		log.Panic("translator not found")
	}

	_ = enTranslation.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validate{
		instance:    validate,
		langDefault: langDefault,
		translator:  trans,
	}
}

func (v *Validate) Struct(payloadStruct interface{}) *status.Status {
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
		return v.ErrorFormatter(err)
	}

	return nil
}

func (v *Validate) ErrorFormatter(err error) *status.Status {
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
		value := e.Translate(v.translator)
		value = regexp.MustCompile(e.Field()).ReplaceAllString(value, field)

		switch e.Tag() {
		case "datetime":
			{
				value = regexDate.ReplaceAllString(value, "yyyy-MM-dd")
				value = regexTime.ReplaceAllString(value, "HH:mm:ss")
			}
		case "required_with":
			{
				value = field + " is a required field"
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
