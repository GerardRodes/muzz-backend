package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/labstack/echo/v4"
)

type profileSvc interface {
	ListPotentialMatches(ctx context.Context, userID uint32) ([]domain.User, error)
}

type profileCtrl struct {
	svc profileSvc
}

func (c profileCtrl) register(g *echo.Group) {
	g.GET("/profiles", c.profiles)
}

func (c profileCtrl) profiles(e echo.Context) error {
	var userID uint32

	{ // parse userID
		userIDStr := e.QueryParam("userID")
		if userIDStr == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "missing query param userID")
		}

		userID64, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid userID")
		}
		userID = uint32(userID64)
	}

	users, err := c.svc.ListPotentialMatches(e.Request().Context(), userID)
	if err != nil {
		return fmt.Errorf("cannot list potential matches: %w", err)
	}

	type resp struct {
		Results []domain.User `json:"results"`
	}
	return e.JSON(http.StatusOK, resp{Results: users})
}
