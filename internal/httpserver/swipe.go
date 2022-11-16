package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c Controller) swipe(e echo.Context) error {
	type req struct {
		UserID     uint32 `json:"userID"`
		ProfileID  uint32 `json:"profileID"`
		Preference bool   `json:"preference"`
	}
	var r req
	if err := e.Bind(&r); err != nil {
		return err
	}

	matchID, err := c.svc.Swipe(e.Request().Context(), r.UserID, r.ProfileID, r.Preference)
	if err != nil {
		return err
	}

	type results struct {
		Matched bool   `json:"matched"`
		MatchID uint64 `json:"matchID,omitempty"`
	}
	type resp struct {
		Results results `json:"results"`
	}
	return e.JSON(http.StatusOK, resp{
		Results: results{
			Matched: matchID > 0,
			MatchID: matchID,
		},
	})
}
