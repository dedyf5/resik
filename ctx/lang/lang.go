// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package lang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dedyf5/resik/ctx/status"
	commonEntity "github.com/dedyf5/resik/entities/common"
	"github.com/dedyf5/resik/pkg/array"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var LangAvailable []language.Tag = []language.Tag{language.English, language.Indonesian}

type langKey string

const (
	ContextKey langKey = "lang"
)

var termDir string = "static/term"

type Lang struct {
	Bundle      *i18n.Bundle
	Localizer   *i18n.Localizer
	LangDefault language.Tag
	LangReq     *language.Tag
	LangAccept  string
}

func NewLang(langDefault language.Tag, langReq *language.Tag, langAccept string) *Lang {
	bundle := i18n.NewBundle(langDefault)
	return &Lang{
		Bundle:     bundle,
		Localizer:  NewLocalizer(bundle, langDefault, langReq, langAccept),
		LangReq:    langReq,
		LangAccept: langAccept,
	}
}

func NewLocalizer(bundle *i18n.Bundle, langDefault language.Tag, langReq *language.Tag, langAccept string) *i18n.Localizer {
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, v := range LangAvailable {
		sourcePath := fmt.Sprintf("%s/%s.json", termDir, v.String())
		bundle.LoadMessageFile(sourcePath)
	}
	return i18n.NewLocalizer(bundle, GetLanguageReqOrDefault(langDefault, langReq).String(), langAccept)
}

// return *Lang, if parsing failed return *status.Status error
//
// error status code: 500
func FromContext(ctx context.Context) (*Lang, *status.Status) {
	if lang, ok := ctx.Value(ContextKey).(*Lang); ok {
		return lang, nil
	}
	return nil, &status.Status{
		Code:       http.StatusInternalServerError,
		CauseError: errors.New("failed to casting lang from context"),
	}
}

func (src *Lang) GetByMessageID(id string) string {
	return src.Localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: id,
	})
}

func (src *Lang) GetByTemplateData(id string, templateData commonEntity.Map) string {
	return src.Localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: templateData,
	})
}

func (src *Lang) LangReqOrDefault() language.Tag {
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
		return language.English
	}
	return *result
}

// return *language.Tag. if parsing failed return *status.Status error
//
// error status code: 400, 500
func GetLanguageAvailable(lang string) (*language.Tag, *status.Status) {
	if _, err := LanguageIsAvailable(lang); err != nil {
		return nil, err
	}
	res, err := language.Parse(lang)
	if err != nil {
		return nil, &status.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return &res, nil
}

// return bool. if lang not avaliable return *status.Status error
//
// error status code: 400
func LanguageIsAvailable(lang string) (bool, *status.Status) {
	if lang == "" {
		msg := "lang is required"
		return false, &status.Status{
			Code:    http.StatusBadRequest,
			Message: msg,
			Detail: map[string]string{
				ContextKey.String(): msg,
			},
		}
	}
	langCodes := make([]string, 0, cap(LangAvailable))
	for _, v := range LangAvailable {
		langCodes = append(langCodes, v.String())
	}
	if array.InArray(lang, langCodes) < 0 {
		msg := fmt.Sprintf("lang must be one of [%s]", strings.Join(langCodes, ", "))
		return false, &status.Status{
			Code:    http.StatusBadRequest,
			Message: msg,
			Detail: map[string]string{
				ContextKey.String(): msg,
			},
		}
	}
	return true, nil
}
