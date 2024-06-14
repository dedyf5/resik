package common

type Request struct {
	Lang string `json:"lang" query:"lang" validate:"omitempty,oneof=en id ja"`
}
