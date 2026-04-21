// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package lang

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dedyf5/resik/ctx/lang/term"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Default = language.English
var Available []language.Tag = []language.Tag{Default, language.Indonesian, language.Japanese}

type langKey string

const (
	ContextKey            langKey = "lang"
	LocaleDir             string  = "static/locales"
	ValidationPrefix      string  = "validation."
	ValidationFieldPrefix string  = ValidationPrefix + "field."
)

type Lang struct {
	Bundle      *i18n.Bundle
	Localizer   *i18n.Localizer
	LangDefault language.Tag
	LangReq     *language.Tag
	LangAccept  string
}

func NewLang(langDefault language.Tag, langReq *language.Tag, langAccept string) *Lang {
	return NewLangLocaleDir(langDefault, langReq, langAccept, LocaleDir)
}

func NewLangLocaleDir(langDefault language.Tag, langReq *language.Tag, langAccept string, localeDir string) *Lang {
	bundle := i18n.NewBundle(langDefault)
	return &Lang{
		Bundle:      bundle,
		Localizer:   NewLocalizer(bundle, langDefault, langReq, langAccept, localeDir),
		LangDefault: langDefault,
		LangReq:     langReq,
		LangAccept:  langAccept,
	}
}

func NewLocalizer(bundle *i18n.Bundle, langDefault language.Tag, langReq *language.Tag, langAccept string, localeDir string) *i18n.Localizer {
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, v := range Available {
		sourcePath := fmt.Sprintf("%s/active.%s.toml", localeDir, v.String())
		if _, err := bundle.LoadMessageFile(sourcePath); err != nil {
			log.Printf("[lang] failed to load message file %s: %v", sourcePath, err)
		}
	}
	return i18n.NewLocalizer(bundle, GetLanguageReqOrDefault(langDefault, langReq).String(), langAccept)
}

// return *Lang, if parsing failed return *response.Status error
//
// error status code: 500
func FromContext(ctx context.Context) (*Lang, *resPkg.Status) {
	if lang, ok := ctx.Value(ContextKey).(*Lang); ok {
		return lang, nil
	}
	return nil, resPkg.NewStatusError(
		http.StatusInternalServerError,
		errors.New("failed to casting lang from context"),
	)
}

func (src *Lang) GetByMessageID(id string) string {
	return term.GetByMessageID(src.Localizer, id)
}

func (src *Lang) GetByTemplateData(id string, templateData any) string {
	return term.GetByTemplateData(src.Localizer, id, templateData)
}

func (src *Lang) ValidationField(field string) string {
	return ValidationFieldPrefix + field
}

func (src *Lang) GetValidationFieldName(field string) string {
	return src.GetByMessageID(src.ValidationField(field))
}

func (src *Lang) GetValidationFieldNameWithQuote(field string) string {
	return term.ValidationQuoteVal.Localize(src.Localizer, src.GetValidationFieldName(field))
}

func (src *Lang) LanguageReqOrDefault() language.Tag {
	return GetLanguageReqOrDefault(src.LangDefault, src.LangReq)
}

func (k langKey) String() string {
	return string(k)
}

func GetLanguageReqOrDefault(langDefault language.Tag, langReq *language.Tag) language.Tag {
	if langReq != nil {
		return *langReq
	}
	return langDefault
}

func GetLanguageOrDefault(lang string) language.Tag {
	result, err := GetLanguageAvailable(lang)
	if err != nil {
		return Default
	}
	return *result
}

// return *language.Tag. if parsing failed return *response.Status error
//
// error status code: 400, 500
func GetLanguageAvailable(lang string) (*language.Tag, *resPkg.Status) {
	if _, err := LanguageIsAvailable(lang); err != nil {
		return nil, err
	}
	res, err := language.Parse(lang)
	if err != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, err)
	}
	return &res, nil
}

// return bool. if lang not avaliable return *response.Status error
//
// error status code: 400
func LanguageIsAvailable(lang string) (bool, *resPkg.Status) {
	if lang == "" {
		msg := "\"Language\" is required"
		technicalErr := errors.New(msg)
		return false, resPkg.NewStatusBadRequest(
			Default.String(),
			ContextKey.String(),
			msg,
			"REQUIRED",
			technicalErr,
		)
	}
	langCodes := make([]string, 0, cap(Available))
	for _, v := range Available {
		langCodes = append(langCodes, v.String())
	}
	if !slices.Contains(langCodes, lang) {
		msg := fmt.Sprintf("\"Language\" must be one of [%s]", strings.Join(langCodes, ", "))
		technicalErr := errors.New(msg)
		return false, resPkg.NewStatusBadRequest(
			Default.String(),
			ContextKey.String(),
			msg,
			"OUT_OF_RANGE",
			technicalErr,
		)
	}
	return true, nil
}
