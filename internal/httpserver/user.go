package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/labstack/echo/v4"
)

type userSvc interface {
	CreateRandom(ctx context.Context) (domain.User, string, error)
}

type userCtrl struct {
	svc userSvc
}

func (c userCtrl) register(g *echo.Group) {
	user := g.Group("/user")
	user.POST("/create", c.createRandomUser)
}

func (c userCtrl) createRandomUser(e echo.Context) error {
	user, password, err := c.svc.CreateRandom(e.Request().Context())
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
