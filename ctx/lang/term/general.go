// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package term

import (
	"github.com/dedyf5/resik/entities/common"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	homeMessageID = "home_message"
	// HomeMessage defines the parameterized greeting/status message for application startup or health endpoints.
	HomeMessage = &struct {
		Message  msgFunc
		Localize func(localizer *i18n.Localizer, appName, appVersion, moduleName, moduleType string) string
	}{
		Message: func() *i18n.Message {
			return &i18n.Message{
				ID:    homeMessageID,
				Other: "{{.AppName}} Version {{.AppVersion}} :: Module {{.ModuleName}} Type {{.ModuleType}}",
			}
		},
		Localize: func(localizer *i18n.Localizer, appName, appVersion, moduleName, moduleType string) string {
			return GetByTemplateData(localizer,
				homeMessageID,
				common.Map{
					"AppName":    appName,
					"AppVersion": appVersion,
					"ModuleName": moduleName,
					"ModuleType": moduleType,
				},
			)
		},
	}

	// IncorrectUsernameOrPassword represents an authentication error term.
	IncorrectUsernameOrPassword = &Term{
		Message: &i18n.Message{
			ID:    "incorrect_username_or_password",
			Other: "Incorrect username or password",
		},
	}

	// Unauthorized represents a general unauthorized access term.
	Unauthorized = &Term{
		Message: &i18n.Message{
			ID:    "unauthorized",
			Other: "Unauthorized",
		},
	}

	// AccessDenied represents an explicit privilege denial term.
	AccessDenied = &Term{
		Message: &i18n.Message{
			ID:    "access_denied",
			Other: "Access denied",
		},
	}

	// InvalidID represents a malformed or invalid identifier error term.
	InvalidID = &Term{
		Message: &i18n.Message{
			ID:    "invalid_id",
			Other: "Invalid ID",
		},
	}

	// InvalidOrExpiredSessionLoginAgain represents a session timeout or invalid token term.
	InvalidOrExpiredSessionLoginAgain = &Term{
		Message: &i18n.Message{
			ID:    "invalid_or_expired_session_login_again",
			Other: "Invalid or expired session, please login again",
		},
	}

	// InvalidTimeFormat represents a time parsing failure term.
	InvalidTimeFormat = &Term{
		Message: &i18n.Message{
			ID:    "invalid_time_format",
			Other: "Invalid time format",
		},
	}

	// TooManyRequests represents a rate limit exceeded error term.
	TooManyRequests = &Term{
		Message: &i18n.Message{
			ID:    "too_many_requests",
			Other: "Too many requests",
		},
	}

	// Organization represents the organization entity label term.
	Organization = &Term{
		Message: &i18n.Message{
			ID:    "organization",
			Other: "Organization",
		},
	}

	// User represents the user entity label term.
	User = &Term{
		Message: &i18n.Message{
			ID:    "user",
			Other: "User",
		},
	}

	// Transaction represents the transaction entity label term.
	Transaction = &Term{
		Message: &i18n.Message{
			ID:    "transaction",
			Other: "Transaction",
		},
	}

	notFoundValID = "not_found_val"
	// NotFoundVal generates a localized entity not found message with a dynamic value.
	NotFoundVal = &struct {
		Message  msgFunc
		Localize locFunc
	}{
		Message: func() *i18n.Message {
			return &i18n.Message{
				ID:    notFoundValID,
				Other: "{{." + keyVal + "}} not found",
			}
		},
		Localize: func(localizer *i18n.Localizer, val string) string {
			return GetByTemplateData(localizer, notFoundValID, common.Map{keyVal: val})
		},
	}

	successfullyCreatedValID = "successfully_created_val"
	// SuccessfullyCreatedVal generates a localized success creation message for an entity.
	SuccessfullyCreatedVal = &struct {
		Message  msgFunc
		Localize locFunc
	}{
		Message: func() *i18n.Message {
			return &i18n.Message{
				ID:    successfullyCreatedValID,
				Other: "Successfully created {{." + keyVal + "}}",
			}
		},
		Localize: func(localizer *i18n.Localizer, val string) string {
			return GetByTemplateData(localizer, successfullyCreatedValID, common.Map{keyVal: val})
		},
	}

	successfullyUpdatedValID = "successfully_updated_val"
	// SuccessfullyUpdatedVal generates a localized success update message for an entity.
	SuccessfullyUpdatedVal = &struct {
		Message  msgFunc
		Localize locFunc
	}{
		Message: func() *i18n.Message {
			return &i18n.Message{
				ID:    successfullyUpdatedValID,
				Other: "Successfully updated {{." + keyVal + "}}",
			}
		},
		Localize: func(localizer *i18n.Localizer, val string) string {
			return GetByTemplateData(localizer, successfullyUpdatedValID, common.Map{keyVal: val})
		},
	}
)
