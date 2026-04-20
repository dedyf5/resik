// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package term

import (
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Term struct {
	Message *i18n.Message
}

func (t *Term) Localize(localizer *i18n.Localizer) string {
	return GetByMessageID(localizer, t.Message.ID)
}

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
