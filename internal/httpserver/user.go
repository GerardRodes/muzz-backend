package httpserver

import (
	"fmt"
	"net/http"

	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/labstack/echo/v4"
)

func (c Controller) createRandomUser(e echo.Context) error {
	user, password, err := c.svc.CreateRandomUser(e.Request().Context())
	if err != nil {
		return fmt.Errorf("creating random user: %w", err)
	}

	type result struct {
		domain.User
		Password string `json:"password"`
	}
	type resp struct {
		Result result `json:"result"`
	}
	return e.JSON(http.StatusOK, resp{
		Result: result{
			User:     user,
			Password: password,
		},
	})
}
