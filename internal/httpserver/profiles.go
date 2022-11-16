package httpserver

import (
	"fmt"
	"net/http"

	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/labstack/echo/v4"
)

func (c Controller) profiles(e echo.Context) error {
	userID := domain.UserIDFromContext(e.Request().Context())

	users, err := c.svc.ListPotentialMatches(e.Request().Context(), userID)
	if err != nil {
		return fmt.Errorf("cannot list potential matches: %w", err)
	}

	type resp struct {
		Results []domain.User `json:"results"`
	}
	return e.JSON(http.StatusOK, resp{Results: users})
}
