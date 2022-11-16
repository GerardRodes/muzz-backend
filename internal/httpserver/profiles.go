package httpserver

import (
	"fmt"
	"net/http"

	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/labstack/echo/v4"
)

func (c Controller) profiles(e echo.Context) error {
	type req struct {
		AgeMin uint8         `query:"ageMin"`
		AgeMax uint8         `query:"ageMax"`
		Gender domain.Gender `query:"gender"`
	}
	var r req
	if err := e.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if r.Gender != "" {
		if err := r.Gender.Validate(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	if r.AgeMax != 0 && r.AgeMin != 0 && r.AgeMax < r.AgeMin {
		return echo.NewHTTPError(http.StatusBadRequest, "max age cannot be lower than min age")
	}

	userID := domain.UserIDFromContext(e.Request().Context())

	users, err := c.svc.ListPotentialMatches(e.Request().Context(), userID, domain.ListPotentialMatchesFilter{
		AgeMin: r.AgeMin,
		AgeMax: r.AgeMax,
		Gender: r.Gender,
	})
	if err != nil {
		return fmt.Errorf("cannot list potential matches: %w", err)
	}

	type resp struct {
		Results []domain.User `json:"results"`
	}
	return e.JSON(http.StatusOK, resp{Results: users})
}
