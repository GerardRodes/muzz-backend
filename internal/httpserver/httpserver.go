package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Config struct {
	HTTPPort string
	UserSvc  UserSvc
}

func Init(c Config) error {
	e := echo.New()
	base := e.Group("")

	userController{c.UserSvc}.register(base)

	if err := e.Start(":" + c.HTTPPort); err != nil {
		return fmt.Errorf("cannot start echo server: %w", err)
	}

	return nil
}
