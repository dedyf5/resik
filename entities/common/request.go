package common

import (
	"reflect"
	"strings"
)

type Request struct {
	Lang string `json:"lang" query:"lang" validate:"omitempty,oneof=en id ja"`
}

func (r *Request) LangAvailable() []string {
	typ := reflect.TypeOf(*r)
	fld := typ.Field(0)
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
