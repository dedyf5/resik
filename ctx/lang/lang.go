// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package lang provides internationalization management, localization wrappers,
// context handlers, and language validation utilities.
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

var (
	// Default represents the default fallback language tag (English).
	Default = language.English
	// Available defines the list of supported language tags in the application.
	Available []language.Tag = []language.Tag{Default, language.Indonesian, language.Japanese}
)

// langKey represents the custom type for context keys related to localization.
type langKey string

const (
	// langKey represents the custom type for context keys related to localization.
	ContextKey langKey = "lang"
	// LocaleDir specifies the directory path containing localization translation files.
	LocaleDir string = "static/locales"
	// ValidationPrefix is the base prefix string for validation translation keys.
	ValidationPrefix string = "validation."
	// ValidationFieldPrefix is the translation key prefix specific to validation fields.
	ValidationFieldPrefix string = ValidationPrefix + "field."
)

// Lang encapsulates the i18n bundle, localizer, and language configurations for a request context.
type Lang struct {
	Bundle      *i18n.Bundle
	Localizer   *i18n.Localizer
	LangDefault language.Tag
	LangReq     *language.Tag
	LangAccept  string
}

// NewLang creates a new Lang instance using the default locale directory.
func NewLang(langDefault language.Tag, langReq *language.Tag, langAccept string) *Lang {
	return NewLangLocaleDir(langDefault, langReq, langAccept, LocaleDir)
}

// NewLangLocaleDir creates a new Lang instance with a custom locale files directory.
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

// NewLocalizer initializes and returns an i18n Localizer by loading translation TOML files from the specified directory.
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

// FromContext retrieves the *Lang instance from the provided context.
// Returns a 500 Internal Server Error status if casting fails.
func FromContext(ctx context.Context) (*Lang, *resPkg.Status) {
	if lang, ok := ctx.Value(ContextKey).(*Lang); ok {
		return lang, nil
	}
	return nil, resPkg.NewStatusError(
		http.StatusInternalServerError,
		errors.New("failed to casting lang from context"),
	)
}

// GetByMessageID retrieves the localized string for a given message ID using the internal localizer.
func (src *Lang) GetByMessageID(id string) string {
	return term.GetByMessageID(src.Localizer, id)
}

// GetByTemplateData retrieves the localized string using dynamic template data and the internal localizer.
func (src *Lang) GetByTemplateData(id string, templateData any) string {
	return term.GetByTemplateData(src.Localizer, id, templateData)
}

// ValidationField returns the fully qualified validation field translation key prefix.
func (src *Lang) ValidationField(field string) string {
	return ValidationFieldPrefix + field
}

// GetValidationFieldName retrieves the localized field name for validation error reporting.
func (src *Lang) GetValidationFieldName(field string) string {
	return src.GetByMessageID(src.ValidationField(field))
}

// GetValidationFieldNameWithQuote retrieves the localized and quoted validation field name.
func (src *Lang) GetValidationFieldNameWithQuote(field string) string {
	return term.ValidationQuoteVal.Localize(src.Localizer, src.GetValidationFieldName(field))
}

// LanguageReqOrDefault returns the requested language tag if present, otherwise returns the default language tag.
func (src *Lang) LanguageReqOrDefault() language.Tag {
	return GetLanguageReqOrDefault(src.LangDefault, src.LangReq)
}

// String converts the langKey to its string representation.
func (k langKey) String() string {
	return string(k)
}

// GetLanguageReqOrDefault returns the requested language tag if non-nil, otherwise falls back to the default language.
func GetLanguageReqOrDefault(langDefault language.Tag, langReq *language.Tag) language.Tag {
	if langReq != nil {
		return *langReq
	}
	return langDefault
}

// GetLanguageOrDefault parses a language string, returning the matched language tag or the default tag if parsing fails.
func GetLanguageOrDefault(lang string) language.Tag {
	result, err := GetLanguageAvailable(lang)
	if err != nil {
		return Default
	}
	return *result
}

// GetLanguageAvailable validates and parses a language string into a language.Tag.
// Returns a 400 Bad Request status if unavailable, or a 500 Internal Server Error if parsing fails.
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

// LanguageIsAvailable checks whether a given language string is provided and supported within the Available languages list.
// Returns a 400 Bad Request status error if invalid or unsupported.
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
