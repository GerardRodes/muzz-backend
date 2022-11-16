package httpserver

import (
	"fmt"
	"net/http"

	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/labstack/echo/v4"
)

func (c Controller) profiles(e echo.Context) error {
	type req struct {
		UserID uint32 `query:"userID"`
	}
	var r req
	if err := e.Bind(&r); err != nil {
		return err
	}

	users, err := c.svc.ListPotentialMatches(e.Request().Context(), r.UserID)
	if err != nil {
		return fmt.Errorf("cannot list potential matches: %w", err)
	}

	type resp struct {
		Results []domain.User `json:"results"`
	}
	return e.JSON(http.StatusOK, resp{Results: users})
}
