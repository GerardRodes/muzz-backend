package httpserver

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Config struct {
	HTTPPort   string
	UserSvc    userSvc
	ProfileSvc profileSvc
}

func Init(c Config) error {
	if c.UserSvc == nil {
		return errors.New("missing user service")
	}
	if c.ProfileSvc == nil {
		return errors.New("missing profile service")
	}

	e := echo.New()
	base := e.Group("")

	userCtrl{c.UserSvc}.register(base)
	profileCtrl{c.ProfileSvc}.register(base)

	if err := e.Start(":" + c.HTTPPort); err != nil {
		return fmt.Errorf("cannot start echo server: %w", err)
	}

	return nil
}
