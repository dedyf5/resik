// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

type LoginPost struct {
	Username string `json:"username" validate:"required" example:"sakisaki"`
	Password string `json:"password" validate:"required" example:"secret"`
}
