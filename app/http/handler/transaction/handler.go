// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/core/transaction/service"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service service.IService
	config  config.Config
}

func New(service service.IService, config config.Config) *Handler {
	return &Handler{
		service: service,
		config:  config,
	}
}

func (h *Handler) GetMerchantOmzet(ctx echo.Context) error {
	return errors.New("MASUK GetMerchantOmzet")
}

func (h *Handler) GetOutletOmzet(ctx echo.Context) error {
	return errors.New("MASUK GetOutletOmzet")
}

func (h *Handler) Login(ctx echo.Context) error {
	return errors.New("MASUK Login")
}

func basicValidation(request *http.Request) (ID int, page int, err error) {
	path := strings.Split(request.URL.Path, "/")
	ID, err1 := strconv.Atoi(path[2])
	if err1 != nil {
		return 0, 0, errors.New("the ID has an invalid format")
	}
	page, err2 := strconv.Atoi(path[4])
	if err2 != nil {
		return 0, 0, errors.New("the page number has an invalid format")
	}
	if page == 0 {
		err3 := errors.New("the page number must be greater than zero")
		return 0, 0, err3
	}
	return ID, page, nil
}

func getDate() *time.Time {
	layout := "2006-01-02"
	str := "2021-11-01"
	t, _ := time.Parse(layout, str)
	return &t
}
