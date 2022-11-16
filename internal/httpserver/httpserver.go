package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	HTTPPort        string
	Service         Service
	HandlersTimeout time.Duration
}

func Init(c Config) error {
	if c.Service == nil {
		return errors.New("missing service")
	}

	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)
		e.DefaultHTTPErrorHandler(err, c)
	}

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: c.HandlersTimeout,
	}))

	ctrl := Controller{c.Service}
	e.POST("/user/create", ctrl.createRandomUser)
	e.POST("/login", ctrl.login)

	auth := e.Group("", newSessionMiddleware(c.Service))
	auth.GET("/profiles", ctrl.profiles)
	auth.POST("/swipe", ctrl.swipe)

	if err := e.Start(":" + c.HTTPPort); err != nil {
		return fmt.Errorf("cannot start echo server: %w", err)
	}

	return nil
}

func newSessionMiddleware(svc Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) (outErr error) {
			cookie, err := e.Cookie("session-id")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			userID, err := svc.LoadSession(e.Request().Context(), cookie.Value)
			if err != nil {
				// remove cookie if present but something has gone wrong
				e.SetCookie(&http.Cookie{
					Name:   "session-id",
					MaxAge: -1,
				})
				if errors.Is(err, domain.ErrNotFound) {
					return echo.NewHTTPError(http.StatusUnauthorized)
				}
				return err
			}
			if userID == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			ctx := domain.ContextWithUserID(e.Request().Context(), userID)
			e.SetRequest(e.Request().WithContext(ctx))
			return next(e)
		}
	}
}
