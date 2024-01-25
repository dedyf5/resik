// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package lang

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	commonEntity "github.com/dedyf5/resik/entities/common"
	statusEntity "github.com/dedyf5/resik/entities/status"
	"github.com/dedyf5/resik/utils/array"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var LangAvailable []language.Tag = []language.Tag{language.English, language.Indonesian}

var termDir string = "static/term"

type Lang struct {
	Bundle     *i18n.Bundle
	Localizer  *i18n.Localizer
	LangReq    language.Tag
	LangAccept string
}

func NewLang(langDefault, langReq language.Tag, langAccept string) *Lang {
	bundle := i18n.NewBundle(langDefault)
	return &Lang{
		Bundle:     bundle,
		Localizer:  NewLocalizer(bundle, langReq, langAccept),
		LangReq:    langReq,
		LangAccept: langAccept,
	}
}

func NewLocalizer(bundle *i18n.Bundle, langReq language.Tag, langAccept string) *i18n.Localizer {
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, v := range LangAvailable {
		sourcePath := fmt.Sprintf("%s/%s.json", termDir, v.String())
		bundle.LoadMessageFile(sourcePath)
	}
	return i18n.NewLocalizer(bundle, langReq.String(), langAccept)
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
		return false, &statusEntity.HTTP{
			Code:    http.StatusBadRequest,
			Message: "lang is required",
		}
	}
	langCodes := make([]string, 0, cap(LangAvailable))
	for _, v := range LangAvailable {
		langCodes = append(langCodes, v.String())
	}
	if array.InArray(lang, langCodes) < 0 {
		msg := fmt.Sprintf("lang must be one of [%s]", strings.Join(langCodes, ", "))
		return false, &statusEntity.HTTP{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}
	return true, nil
}
