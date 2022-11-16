package httpserver

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/labstack/echo/v4"
)

func (c Controller) login(e echo.Context) error {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var r req
	if err := e.Bind(&r); err != nil {
		return err
	}

	sessionID, err := c.svc.Login(e.Request().Context(), r.Email, r.Password)
	if err != nil {
		if errors.Is(err, domain.ErrWrongPassword) {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return fmt.Errorf("login: %w", err)
	}

	e.SetCookie(&http.Cookie{
		Name:     "session-id",
		Value:    sessionID,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // only our site but also when comming from 3rd party sites
		// Secure: true, this should be set when using https
		// Domain: "muzz.com", this should be set accordingly
	})

	return e.NoContent(http.StatusOK)
}
