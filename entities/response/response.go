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
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type Log struct {
	Response Response `json:"response"`
	Message  string   `json:"message"`
}
