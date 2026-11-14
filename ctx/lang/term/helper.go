// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package term

import (
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const keyVal = "Val"

// Term represents an internationalizable message configuration item.
type Term struct {
	Message *i18n.Message
}

// Localize translates the Term message using the provided i18n localizer.
func (t *Term) Localize(localizer *i18n.Localizer) string {
	return GetByMessageID(localizer, t.Message.ID)
}

// GetByMessageID retrieves the localized message string for the given message ID.
// Logs a warning and returns the ID itself if localization fails.
func GetByMessageID(localizer *i18n.Localizer, id string) string {
	val, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: id,
	})
	if err != nil {
		log.Printf("[lang] failed to localize message id %s: %v", id, err)
		return id
	}
	return val
}

// GetByTemplateData retrieves the localized message string using dynamic template data.
// Logs a warning and returns the ID itself if localization fails.
func GetByTemplateData(localizer *i18n.Localizer, id string, templateData any) string {
	val, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: templateData,
	})
	if err != nil {
		log.Printf("[lang] failed to localize message id %s: %v", id, err)
		return id
	}
	return val
}
