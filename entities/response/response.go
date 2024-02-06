package response

type Response struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Meta   *Meta       `json:"meta,omitempty"`
}

type Status struct {
	Code    string            `json:"code" example:"200.1"`
	Message string            `json:"message" example:"OK"`
	Detail  map[string]string `json:"detail,omitempty"`
}

type Meta struct {
	Total uint64 `json:"total" example:"100"`
	Limit int    `json:"limit" example:"10"`
	Page  *Page  `json:"page"`
}

type Page struct {
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
