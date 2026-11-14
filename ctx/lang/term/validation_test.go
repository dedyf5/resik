// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package term

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func getValidationTestLocalizer(t *testing.T) *i18n.Localizer {
	t.Helper()
	bundle := i18n.NewBundle(language.English)

	err := bundle.AddMessages(language.English,
		&i18n.Message{ID: "validation.field.page", Other: "Page"},
		&i18n.Message{ID: "validation.field.limit", Other: "Limit"},
		&i18n.Message{ID: "validation.field.order", Other: "Order"},
		&i18n.Message{ID: "validation.field.id", Other: "ID"},
		&i18n.Message{ID: "validation.field.name", Other: "Name"},
		&i18n.Message{ID: "validation.field.description", Other: "Description"},
		&i18n.Message{ID: "validation.field.created_at", Other: "Created at"},
		&i18n.Message{ID: "validation.field.updated_at", Other: "Updated at"},
		&i18n.Message{ID: "validation.field.mode", Other: "Mode"},
		&i18n.Message{ID: "validation.field.datetime_start", Other: "Start at"},
		&i18n.Message{ID: "validation.field.datetime_end", Other: "End at"},
		&i18n.Message{ID: validationFieldTimezoneID, Other: "Timezone"},
		&i18n.Message{ID: validationFieldBranchIDID, Other: "Branch ID"},
		&i18n.Message{ID: validationFieldOrganizationIDID, Other: "Organization ID"},
		&i18n.Message{ID: validationFieldBranchBranchIDID, Other: "Branch ID"},
		&i18n.Message{ID: validationFieldBranchBranchNameID, Other: "Branch Name"},
		&i18n.Message{ID: validationFieldBranchOrganizationIDID, Other: "Organization ID"},
		&i18n.Message{ID: validationFieldBranchOrganizationNameID, Other: "Organization Name"},
		&i18n.Message{ID: validationFieldBillTotalID, Other: "Bill Total"},
		&i18n.Message{ID: validationFieldTransactedAtID, Other: "Transacted at"},
		&i18n.Message{ID: validationFieldUsernameID, Other: "Username"},
		&i18n.Message{ID: validationFieldPasswordID, Other: "Password"},
		&i18n.Message{ID: validationQuoteValID, Other: "\"{{.Val}}\""},
		&i18n.Message{ID: validationTypeMessageID, Other: "{{.Field}} type must be {{.Expected}}, got {{.Actual}}"},
	)
	if err != nil {
		t.Fatalf("failed to add messages: %v", err)
	}

	return i18n.NewLocalizer(bundle, "en")
}

func TestValidationTermsTable(t *testing.T) {
	localizer := getValidationTestLocalizer(t)

	validationTerms := map[string]*Term{
		"ValidationFieldPage":                   ValidationFieldPage,
		"ValidationFieldLimit":                  ValidationFieldLimit,
		"ValidationFieldOrder":                  ValidationFieldOrder,
		"ValidationFieldID":                     ValidationFieldID,
		"ValidationFieldName":                   ValidationFieldName,
		"ValidationFieldDescription":            ValidationFieldDescription,
		"ValidationFieldCreatedAt":              ValidationFieldCreatedAt,
		"ValidationFieldUpdatedAt":              ValidationFieldUpdatedAt,
		"ValidationFieldMode":                   ValidationFieldMode,
		"ValidationFieldDatetimeStart":          ValidationFieldDatetimeStart,
		"ValidationFieldDatetimeEnd":            ValidationFieldDatetimeEnd,
		"ValidationFieldTimezone":               ValidationFieldTimezone,
		"ValidationFieldBranchID":               ValidationFieldBranchID,
		"ValidationFieldOrganizationID":         ValidationFieldOrganizationID,
		"ValidationFieldBranchBranchID":         ValidationFieldBranchBranchID,
		"ValidationFieldBranchBranchName":       ValidationFieldBranchBranchName,
		"ValidationFieldBranchOrganizationID":   ValidationFieldBranchOrganizationID,
		"ValidationFieldBranchOrganizationName": ValidationFieldBranchOrganizationName,
		"ValidationFieldBillTotal":              ValidationFieldBillTotal,
		"ValidationFieldTransactedAt":           ValidationFieldTransactedAt,
		"ValidationFieldUsername":               ValidationFieldUsername,
		"ValidationFieldPassword":               ValidationFieldPassword,
	}

	for name, term := range validationTerms {
		t.Run(name, func(t *testing.T) {
			if term.Message == nil || term.Message.ID == "" {
				t.Fatalf("%s message ID is empty", name)
			}
			res := term.Localize(localizer)
			if res == "" {
				t.Errorf("expected localized string for %s, got empty", name)
			}
		})
	}

	t.Run("ValidationTypeMessage", func(t *testing.T) {
		msg := ValidationTypeMessage.Message()
		if msg == nil || msg.ID != validationTypeMessageID {
			t.Errorf("expected message ID %s", validationTypeMessageID)
		}
		res := ValidationTypeMessage.Localize(localizer, "Age", "int", "string")
		expected := "Age type must be int, got string"
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
				name:      "ValidationQuoteVal",
				messageID: validationQuoteValID,
				msgFunc:   ValidationQuoteVal.Message,
				locFunc:   ValidationQuoteVal.Localize,
				inputVal:  "test_val",
				expected:  "\"test_val\"",
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
