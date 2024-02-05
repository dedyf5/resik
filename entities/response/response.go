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
	Limit int    `json:"limit"`
	Page  *Page  `json:"page"`
}

type Page struct {
	First    int  `json:"first"`
	Previous *int `json:"previous"`
	Current  int  `json:"current"`
	Next     *int `json:"next"`
	Last     int  `json:"last"`
}

type Log struct {
	Response Response `json:"response"`
	Message  string   `json:"message"`
}
