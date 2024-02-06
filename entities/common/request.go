package common

type Request struct {
	Lang string `json:"lang" query:"lang" validate:"oneof=en id"`
}
