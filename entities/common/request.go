package common

import (
	"reflect"
	"strings"
)

// Request difines common query and payload parameters shared across incoming requests.
type Request struct {
	Lang string `json:"lang" query:"lang" validate:"omitempty,oneof=en id ja"`
}

// LangAvailable extracts the list of allowed language codes from the Lang field's oneof validation tag.
func (r *Request) LangAvailable() []string {
	fld := reflect.TypeFor[Request]().Field(0)
	validate := fld.Tag.Get("validate")
	validateList := strings.Split(validate, ",")
	str := ""
	for _, v1 := range validateList {
		prefix := "oneof="
		if strings.Contains(v1, prefix) {
			str = strings.ReplaceAll(v1, prefix, "")
		}
	}
	return strings.Split(str, " ")
}
