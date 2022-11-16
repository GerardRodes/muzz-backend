package httpserver

import (
	"errors"
	"fmt"
	"time"

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
	e.GET("/profiles", ctrl.profiles)
	e.POST("/swipe", ctrl.swipe)

	if err := e.Start(":" + c.HTTPPort); err != nil {
		return fmt.Errorf("cannot start echo server: %w", err)
	}

	return nil
}
