// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package term

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/dedyf5/resik/entities/common"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	simpleSuccessMessage = "Simple success message"
)

func getTestLocalizer(t *testing.T) *i18n.Localizer {
	t.Helper()
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	err := bundle.AddMessages(language.English, &i18n.Message{
		ID:    "test_message",
		Other: "Hello {{.Val}}",
	})
	if err != nil {
		t.Fatalf("failed to add message: %v", err)
	}

	err = bundle.AddMessages(language.English, &i18n.Message{
		ID:    "simple_message",
		Other: simpleSuccessMessage,
	})
	if err != nil {
		t.Fatalf("failed to add message: %v", err)
	}

	return i18n.NewLocalizer(bundle, "en")
}

func TestTerm_Localize(t *testing.T) {
	localizer := getTestLocalizer(t)

	term := &Term{
		Message: &i18n.Message{
			ID:    "simple_message",
			Other: simpleSuccessMessage,
		},
	}

	result := term.Localize(localizer)
	if result != simpleSuccessMessage {
		t.Errorf("expected %q, got %q", simpleSuccessMessage, result)
	}
}

func TestGetByMessageID(t *testing.T) {
	localizer := getTestLocalizer(t)

	t.Run("Success localization", func(t *testing.T) {
		result := GetByMessageID(localizer, "simple_message")
		if result != simpleSuccessMessage {
			t.Errorf("expected %q, got %q", simpleSuccessMessage, result)
		}
	})

	t.Run("Fallback on missing message ID", func(t *testing.T) {
		missingID := "non_existent_id"
		result := GetByMessageID(localizer, missingID)
		if result != missingID {
			t.Errorf("expected fallback to %q, got %q", missingID, result)
		}
	})
}

func TestGetByTemplateData(t *testing.T) {
	localizer := getTestLocalizer(t)

	t.Run("Success localization with template data", func(t *testing.T) {
		data := common.Map{"Val": "World"}
		result := GetByTemplateData(localizer, "test_message", data)
		expected := "Hello World"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("Fallback on missing message ID with template data", func(t *testing.T) {
		missingID := "non_existent_template_id"
		data := common.Map{"Val": "World"}
		result := GetByTemplateData(localizer, missingID, data)
		if result != missingID {
			t.Errorf("expected fallback to %q, got %q", missingID, result)
		}
	})
}
