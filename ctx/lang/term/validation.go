// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package term

import (
	"github.com/dedyf5/resik/entities/common"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	// ValidationFieldPage represents the pagination page field label.
	ValidationFieldPage = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.page",
			Other: "Page",
		},
	}

	// ValidationFieldLimit represents the pagination limit field label.
	ValidationFieldLimit = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.limit",
			Other: "Limit",
		},
	}

	// ValidationFieldOrder represents the sorting order field label.
	ValidationFieldOrder = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.order",
			Other: "Order",
		},
	}

	// ValidationFieldID represents the generic ID validation field label.
	ValidationFieldID = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.id",
			Other: "ID",
		},
	}

	// ValidationFieldName represents the name validation field label.
	ValidationFieldName = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.name",
			Other: "Name",
		},
	}

	// ValidationFieldDescription represents the description validation field label.
	ValidationFieldDescription = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.description",
			Other: "Description",
		},
	}

	// ValidationFieldCreatedAt represents the created timestamp validation field label.
	ValidationFieldCreatedAt = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.created_at",
			Other: "Created at",
		},
	}

	// ValidationFieldUpdatedAt represents the updated timestamp validation field label.
	ValidationFieldUpdatedAt = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.updated_at",
			Other: "Updated at",
		},
	}

	// ValidationFieldMode represents the operational mode validation field label.
	ValidationFieldMode = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.mode",
			Other: "Mode",
		},
	}

	// ValidationFieldDatetimeStart represents the start datetime validation field label.
	ValidationFieldDatetimeStart = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.datetime_start",
			Other: "Start at",
		},
	}

	// ValidationFieldDatetimeEnd represents the end datetime validation field label.
	ValidationFieldDatetimeEnd = &Term{
		Message: &i18n.Message{
			ID:    "validation.field.datetime_end",
			Other: "End at",
		},
	}

	validationFieldTimezoneID = "validation.field.timezone"
	// ValidationFieldTimezone represents the timezone validation field label.
	ValidationFieldTimezone = &Term{
		Message: &i18n.Message{
			ID:    validationFieldTimezoneID,
			Other: "Timezone",
		},
	}

	branchID = "Branch ID"

	validationFieldBranchIDID = "validation.field.branch_id"
	// ValidationFieldBranchID represents the branch ID validation field label.
	ValidationFieldBranchID = &Term{
		Message: &i18n.Message{
			ID:    validationFieldBranchIDID,
			Other: branchID,
		},
	}

	organizationID = "Organization ID"

	validationFieldOrganizationIDID = "validation.field.organization_id"
	// ValidationFieldOrganizationID represents the organization ID validation field label.
	ValidationFieldOrganizationID = &Term{
		Message: &i18n.Message{
			ID:    validationFieldOrganizationIDID,
			Other: organizationID,
		},
	}

	validationFieldBranchBranchIDID = "validation.field.branch.id"
	// ValidationFieldBranchBranchID represents the nested branch identifier validation field label.
	ValidationFieldBranchBranchID = &Term{
		Message: &i18n.Message{
			ID:    validationFieldBranchBranchIDID,
			Other: branchID,
		},
	}

	validationFieldBranchBranchNameID = "validation.field.branch.name"
	// ValidationFieldBranchBranchName represents the nested branch name validation field label.
	ValidationFieldBranchBranchName = &Term{
		Message: &i18n.Message{
			ID:    validationFieldBranchBranchNameID,
			Other: "Branch Name",
		},
	}

	validationFieldBranchOrganizationIDID = "validation.field.branch.organization.id"
	// ValidationFieldBranchOrganizationID represents the nested organization identifier within a branch validation field label.
	ValidationFieldBranchOrganizationID = &Term{
		Message: &i18n.Message{
			ID:    validationFieldBranchOrganizationIDID,
			Other: organizationID,
		},
	}

	validationFieldBranchOrganizationNameID = "validation.field.branch.organization.name"
	// ValidationFieldBranchOrganizationName represents the nested organization name within a branch validation field label.
	ValidationFieldBranchOrganizationName = &Term{
		Message: &i18n.Message{
			ID:    validationFieldBranchOrganizationNameID,
			Other: "Organization Name",
		},
	}

	validationFieldBillTotalID = "validation.field.bill_total"
	// ValidationFieldBillTotal represents the bill total validation field label.
	ValidationFieldBillTotal = &Term{
		Message: &i18n.Message{
			ID:    validationFieldBillTotalID,
			Other: "Bill Total",
		},
	}

	validationFieldTransactedAtID = "validation.field.transacted_at"
	// ValidationFieldTransactedAt represents the transaction timestamp validation field label.
	ValidationFieldTransactedAt = &Term{
		Message: &i18n.Message{
			ID:    validationFieldTransactedAtID,
			Other: "Transacted at",
		},
	}

	validationFieldUsernameID = "validation.field.username"
	// ValidationFieldUsername represents the username validation field label.
	ValidationFieldUsername = &Term{
		Message: &i18n.Message{
			ID:    validationFieldUsernameID,
			Other: "Username",
		},
	}

	validationFieldPasswordID = "validation.field.password"
	// ValidationFieldPassword represents the password validation field label.
	ValidationFieldPassword = &Term{
		Message: &i18n.Message{
			ID:    validationFieldPasswordID,
			Other: "Password",
		},
	}

	validationQuoteValID = "validation.quote_val"
	// ValidationQuoteVal wraps a validation value in quotation marks for localized error formatting.
	ValidationQuoteVal = &struct {
		Message  msgFunc
		Localize locFunc
	}{
		Message: func() *i18n.Message {
			return &i18n.Message{
				ID:    validationQuoteValID,
				Other: "\"{{." + keyVal + "}}\"",
			}
		},
		Localize: func(localizer *i18n.Localizer, val string) string {
			return GetByTemplateData(localizer, validationQuoteValID, common.Map{keyVal: val})
		},
	}

	validationTypeMessageID = "validation.type_message"
	// ValidationTypeMessage generates a localized type mismatch validation error message.
	ValidationTypeMessage = &struct {
		Message  msgFunc
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
