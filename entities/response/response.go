package response

type Response struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Meta   *Meta       `json:"meta,omitempty"`
}

type Status struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Detail  map[string]string `json:"detail,omitempty"`
}

type Meta struct {
	Total uint64 `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type Log struct {
	Response Response `json:"response"`
	Message  string   `json:"message"`
}
