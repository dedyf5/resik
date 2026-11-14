// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package term

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	organization = "Organization"
	user         = "User"
)

func getGeneralTestLocalizer(t *testing.T) *i18n.Localizer {
	t.Helper()
	bundle := i18n.NewBundle(language.English)

	// Daftarkan message ID yang digunakan di general.go
	err := bundle.AddMessages(language.English,
		&i18n.Message{ID: homeMessageID, Other: "{{.AppName}} Version {{.AppVersion}} :: Module {{.ModuleName}} Type {{.ModuleType}}"},
		&i18n.Message{ID: "incorrect_username_or_password", Other: "Incorrect username or password"},
		&i18n.Message{ID: "unauthorized", Other: "Unauthorized"},
		&i18n.Message{ID: "access_denied", Other: "Access denied"},
		&i18n.Message{ID: "invalid_id", Other: "Invalid ID"},
		&i18n.Message{ID: "invalid_or_expired_session_login_again", Other: "Invalid or expired session, please login again"},
		&i18n.Message{ID: "invalid_time_format", Other: "Invalid time format"},
		&i18n.Message{ID: "too_many_requests", Other: "Too many requests"},
		&i18n.Message{ID: "organization", Other: organization},
		&i18n.Message{ID: "user", Other: user},
		&i18n.Message{ID: "transaction", Other: "Transaction"},
		&i18n.Message{ID: notFoundValID, Other: "{{.Val}} not found"},
		&i18n.Message{ID: successfullyCreatedValID, Other: "Successfully created {{.Val}}"},
		&i18n.Message{ID: successfullyUpdatedValID, Other: "Successfully updated {{.Val}}"},
	)
	if err != nil {
		t.Fatalf("failed to add messages: %v", err)
	}

	return i18n.NewLocalizer(bundle, "en")
}

func TestGeneralTermsAndStructures(t *testing.T) {
	localizer := getGeneralTestLocalizer(t)

	t.Run("Standard Term Localizations", func(t *testing.T) {
		terms := map[string]*Term{
			"IncorrectUsernameOrPassword": IncorrectUsernameOrPassword,
			"Unauthorized":                Unauthorized,
			"AccessDenied":                AccessDenied,
			"InvalidID":                   InvalidID,
			"InvalidOrExpiredSession":     InvalidOrExpiredSessionLoginAgain,
			"InvalidTimeFormat":           InvalidTimeFormat,
			"TooManyRequests":             TooManyRequests,
			"Organization":                Organization,
			"User":                        User,
			"Transaction":                 Transaction,
		}

		for name, term := range terms {
			t.Run(name, func(t *testing.T) {
				if term.Message == nil || term.Message.ID == "" {
					t.Fatalf("%s has empty message ID", name)
				}
				res := term.Localize(localizer)
				if res == "" {
					t.Errorf("expected localized text for %s", name)
				}
			})
		}
	})

	t.Run("HomeMessage Message() and Localize()", func(t *testing.T) {
		msg := HomeMessage.Message()
		if msg == nil || msg.ID != homeMessageID {
			t.Errorf("expected message ID %s, got %v", homeMessageID, msg)
		}

		res := HomeMessage.Localize(localizer, "ResikApp", "v1.0", "Core", "REST")
		expected := "ResikApp Version v1.0 :: Module Core Type REST"
		if res != expected {
			t.Errorf("expected %q, got %q", expected, res)
		}
	})

	t.Run("Parameterized Val Struct Terms", func(t *testing.T) {
		testCases := []struct {
			name      string
			messageID string
			msgFunc   msgFunc
			locFunc   locFunc
			inputVal  string
			expected  string
		}{
			{
				name:      "NotFoundVal",
				messageID: notFoundValID,
				msgFunc:   NotFoundVal.Message,
				locFunc:   NotFoundVal.Localize,
				inputVal:  "Account",
				expected:  "Account not found",
			},
			{
				name:      "SuccessfullyCreatedVal",
				messageID: successfullyCreatedValID,
				msgFunc:   SuccessfullyCreatedVal.Message,
				locFunc:   SuccessfullyCreatedVal.Localize,
				inputVal:  user,
				expected:  "Successfully created User",
			},
			{
				name:      "SuccessfullyUpdatedVal",
				messageID: successfullyUpdatedValID,
				msgFunc:   SuccessfullyUpdatedVal.Message,
				locFunc:   SuccessfullyUpdatedVal.Localize,
				inputVal:  organization,
				expected:  "Successfully updated Organization",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				msg := tc.msgFunc()
				if msg == nil || msg.ID != tc.messageID {
					t.Errorf("expected message ID %s", tc.messageID)
				}

				res := tc.locFunc(localizer, tc.inputVal)
				if res != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, res)
				}
			})
		}
	})
}
