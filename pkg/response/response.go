package response

type Response struct {
	Status ResponseStatus `json:"status"`
	Data   interface{}    `json:"data,omitempty"`
	Meta   *ResponseMeta  `json:"meta,omitempty"`
}

type ResponseStatus struct {
	Code    string            `json:"code" example:"200.1"`
	Message string            `json:"message" example:"OK"`
	Detail  map[string]string `json:"detail,omitempty"`
}

type ResponseMeta struct {
	Total uint64        `json:"total" example:"100"`
	Limit int           `json:"limit" example:"10"`
	Page  *ResponsePage `json:"page"`
}

type ResponsePage struct {
	First    int  `json:"first" example:"1"`
	Previous *int `json:"previous" example:"2"`
	Current  int  `json:"current" example:"3"`
	Next     *int `json:"next" example:"4"`
	Last     int  `json:"last" example:"9"`
}

type Log struct {
	Response Response `json:"response"`
	Message  string   `json:"message"`
}
