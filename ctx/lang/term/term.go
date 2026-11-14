// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package term provides internationalization terms, message structures,
// and localization utility functions for translating error messages, UI elements,
// and validation field names.
package term

import "github.com/nicksnyder/go-i18n/v2/i18n"

type (
	msgFunc func() *i18n.Message
	locFunc func(localizer *i18n.Localizer, val string) string
)
