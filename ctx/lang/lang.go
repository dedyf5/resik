// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package lang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	commonEntity "github.com/dedyf5/resik/entities/common"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Default = language.English
var Available []language.Tag = []language.Tag{Default, language.Indonesian, language.Japanese}

type langKey string

const (
	ContextKey            langKey = "lang"
	TermDir               string  = "static/term"
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
	return NewLangTermDir(langDefault, langReq, langAccept, TermDir)
}

func NewLangTermDir(langDefault language.Tag, langReq *language.Tag, langAccept string, termDir string) *Lang {
	bundle := i18n.NewBundle(langDefault)
	return &Lang{
		Bundle:      bundle,
		Localizer:   NewLocalizer(bundle, langDefault, langReq, langAccept, termDir),
		LangDefault: langDefault,
		LangReq:     langReq,
		LangAccept:  langAccept,
	}
}

func NewLocalizer(bundle *i18n.Bundle, langDefault language.Tag, langReq *language.Tag, langAccept string, termDir string) *i18n.Localizer {
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, v := range Available {
		sourcePath := fmt.Sprintf("%s/%s.json", termDir, v.String())
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
	val, err := src.Localizer.Localize(&i18n.LocalizeConfig{
		MessageID: id,
	})
	if err != nil {
		log.Printf("[lang] failed to localize message id %s: %v", id, err)
		return id
	}
	return val
}

func (src *Lang) GetByTemplateData(id string, templateData commonEntity.Map) string {
	val, err := src.Localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: templateData,
	})
	if err != nil {
		log.Printf("[lang] failed to localize message id %s: %v", id, err)
		return id
	}
	return val
}

func (src *Lang) ValidationField(field string) string {
	return ValidationFieldPrefix + field
}

func (src *Lang) GetValidationFieldName(field string) string {
	return src.GetByMessageID(src.ValidationField(field))
}

func (src *Lang) GetValidationFieldNameWithQuote(field string) string {
	return src.GetByTemplateData(
		src.ValidationField("quote_val"),
		map[string]any{
			"val": src.GetValidationFieldName(field),
		},
	)
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
