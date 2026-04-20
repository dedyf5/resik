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
	HomeMessage   = &struct {
		Message  func() *i18n.Message
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

	IncorrectUsernameOrPassword = &Term{
		Message: &i18n.Message{
			ID:    "incorrect_username_or_password",
			Other: "Incorrect username or password",
		},
	}

	Unauthorized = &Term{
		Message: &i18n.Message{
			ID:    "unauthorized",
			Other: "Unauthorized",
		},
	}

	InvalidOrExpiredSessionLoginAgain = &Term{
		Message: &i18n.Message{
			ID:    "invalid_or_expired_session_login_again",
			Other: "Invalid or expired session, please login again",
		},
	}

	InvalidTimeFormat = &Term{
		Message: &i18n.Message{
			ID:    "invalid_time_format",
			Other: "Invalid time format",
		},
	}

	TooManyRequests = &Term{
		Message: &i18n.Message{
			ID:    "too_many_requests",
			Other: "Too many requests",
		},
	}

	Merchant = &Term{
		Message: &i18n.Message{
			ID:    "merchant",
			Other: "Merchant",
		},
	}

	successfullyCreatedValID = "successfully_created_val"
	SuccessfullyCreatedVal   = &struct {
		Message  func() *i18n.Message
		Localize func(localizer *i18n.Localizer, val string) string
	}{
		Message: func() *i18n.Message {
			return &i18n.Message{
				ID:    successfullyCreatedValID,
				Other: "Successfully created {{.Val}}",
			}
		},
		Localize: func(localizer *i18n.Localizer, val string) string {
			return GetByTemplateData(localizer, successfullyCreatedValID, common.Map{"Val": val})
		},
	}

	successfullyUpdatedValID = "successfully_updated_val"
	SuccessfullyUpdatedVal   = &struct {
		Message  func() *i18n.Message
		Localize func(localizer *i18n.Localizer, val string) string
	}{
		Message: func() *i18n.Message {
			return &i18n.Message{
				ID:    successfullyUpdatedValID,
				Other: "Successfully updated {{.Val}}",
			}
		},
		Localize: func(localizer *i18n.Localizer, val string) string {
			return GetByTemplateData(localizer, successfullyUpdatedValID, common.Map{"Val": val})
		},
	}
)
