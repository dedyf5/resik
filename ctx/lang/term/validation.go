// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package term

import (
	"github.com/dedyf5/resik/entities/common"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	ValidationFieldPage = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.page",
			Other: "Page",
		},
	}

	ValidationFieldLimit = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.limit",
			Other: "Limit",
		},
	}

	ValidationFieldOrder = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.order",
			Other: "Order",
		},
	}

	ValidationFieldID = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.id",
			Other: "ID",
		},
	}

	ValidationFieldName = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.name",
			Other: "Name",
		},
	}

	ValidationFieldDescription = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.description",
			Other: "Description",
		},
	}

	ValidationFieldCreatedAt = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.created_at",
			Other: "Created at",
		},
	}

	ValidationFieldUpdatedAt = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.updated_at",
			Other: "Updated at",
		},
	}

	ValidationFieldMode = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.mode",
			Other: "Mode",
		},
	}

	ValidationFieldDatetimeStart = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.datetime_start",
			Other: "Start at",
		},
	}

	ValidationFieldDatetimeEnd = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.datetime_end",
			Other: "End at",
		},
	}

	validationFieldTimezoneID = "validation.field.timezone"
	ValidationFieldTimezone   = &Term{
		Message: &i18n.Message{
			ID:    validationFieldTimezoneID,
			Other: "Timezone",
		},
	}

	validationFieldMerchantIDID = "validation.field.merchant_id"
	ValidationFieldMerchantID   = &Term{
		Message: &i18n.Message{
			ID:    validationFieldMerchantIDID,
			Other: "Merchant ID",
		},
	}

	validationFieldOutletIDID = "validation.field.outlet_id"
	ValidationFieldOutletID   = &Term{
		Message: &i18n.Message{
			ID:    validationFieldOutletIDID,
			Other: "Outlet ID",
		},
	}

	validationFieldUsernameID = "validation.field.username"
	ValidationFieldUsername   = &Term{
		Message: &i18n.Message{
			ID:    validationFieldUsernameID,
			Other: "Username",
		},
	}

	validationFieldPasswordID = "validation.field.password"
	ValidationFieldPassword   = &Term{
		Message: &i18n.Message{
			ID:    validationFieldPasswordID,
			Other: "Password",
		},
	}

	validationQuoteValID = "validation.quote_val"
	ValidationQuoteVal   = &struct {
		Message  func() *i18n.Message
		Localize func(localizer *i18n.Localizer, val string) string
	}{
		Message: func() *i18n.Message {
			return &i18n.Message{
				ID:    validationQuoteValID,
				Other: "\"{{.Val}}\"",
			}
		},
		Localize: func(localizer *i18n.Localizer, val string) string {
			return GetByTemplateData(localizer, validationQuoteValID, common.Map{"Val": val})
		},
	}

	validationTypeMessageID = "validation.type_message"
	ValidationTypeMessage   = &struct {
		Message  func() *i18n.Message
		Localize func(localizer *i18n.Localizer, field, expected, actual string) string
	}{
		Message: func() *i18n.Message {
			return &i18n.Message{
				ID:    validationTypeMessageID,
				Other: "{{.Field}} type must be {{.Expected}}, got {{.Actual}}",
			}
		},
		Localize: func(localizer *i18n.Localizer, field, expected, actual string) string {
			return GetByTemplateData(
				localizer,
				validationTypeMessageID,
				common.Map{
					"Field":    field,
					"Expected": expected,
					"Actual":   actual,
				},
			)
		},
	}
)
