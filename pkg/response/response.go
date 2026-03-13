// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

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

type Log struct {
	Response Response `json:"response"`
	Message  string   `json:"message"`
}
