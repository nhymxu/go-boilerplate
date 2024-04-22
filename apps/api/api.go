package api

import (
	"net/http"

	"github.com/brpaz/echozap"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/nhymxu/go-boilerplate/pkg/config"
)

func New() *echo.Echo {
	e := echo.New()
	if config.ENV.Debug {
		e.Debug = true
	}

	e.Use(
		middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		middleware.Recover(),
		middleware.RequestID(),
		//middleware.Secure(),
		//middleware.CORS(),
		middleware.Gzip(),
		middleware.RemoveTrailingSlash(),
		echozap.ZapLogger(zap.L()),
		sentryecho.New(sentryecho.Options{}),
	)

	v1 := e.Group("/v1")
	groupV1Routes(v1)

	return e
}
