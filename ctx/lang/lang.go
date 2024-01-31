// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package lang

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

func GetLanguageAvailable(lang string) (*language.Tag, error) {
	if _, err := LanguageIsAvailable(lang); err != nil {
		return nil, err
	}
	res := language.MustParse(lang)
	return &res, nil
}

func LanguageIsAvailable(lang string) (bool, error) {
	if lang == "" {
		return false, errors.New("lang is required")
	}
	langCodes := make([]string, 0, cap(LangAvailable))
	for _, v := range LangAvailable {
		langCodes = append(langCodes, v.String())
	}
	if array.InArray(lang, langCodes) < 0 {
		return false, fmt.Errorf("lang must be one of [%s]", strings.Join(langCodes, ", "))
	}
	return true, nil
}
