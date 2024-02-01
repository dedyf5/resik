// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package translations

import (
	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/language"
)

type translation struct {
	Tag             string
	Translation     string
	Override        bool
	CustomRegisFunc validator.RegisterTranslationsFunc
	CustomTransFunc validator.TranslationFunc
}

var Map map[language.Tag][]translation = map[language.Tag][]translation{
	language.Indonesian: Indonesian,
	language.English:    English,
}

func Translators() (res []locales.Translator) {
	for _, v := range langCtx.Available {
		res = append(res, LanguageToTranslator(v))
	}
	return
}

func LanguageToTranslator(lang language.Tag) locales.Translator {
	switch lang.String() {
	case "id":
		return id.New()
	}
	return en.New()
}
