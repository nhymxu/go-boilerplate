package api

import (
	"fmt"
	"net/http"
	"slices"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/karagenc/fj4echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nhymxu/go-boilerplate/pkg/config"
)

func New() *echo.Echo {
	e := newEchoApp(config.ENV.Debug)

	v1 := e.Group("/v1")
	groupV1Routes(v1)

	return e
}

func newEchoApp(debug bool) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Debug = debug
	e.HideBanner = true

	logger := zap.L()

	e.JSONSerializer = fj4echo.New()

	skipPaths := []string{
		"/favicon.ico",
		"/swagger",
		"/metrics",
		"/health",
		"/ping",
		"/special-endpoint-can-replace-later",
	}

	skipper := func(c echo.Context) bool {
		return slices.Contains(skipPaths, c.Request().URL.Path)
	}

	e.Use(
		middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		middleware.Recover(),
		middleware.RequestID(),
		//middleware.Secure(),
		//middleware.CORS(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Level:   config.APIGzipLevel,
			Skipper: skipper,
		}), middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			Skipper:         skipper,
			LogURI:          true,
			LogStatus:       true,
			LogLatency:      true,
			LogRemoteIP:     true,
			LogMethod:       true,
			LogResponseSize: true,
			LogUserAgent:    true,
			LogRequestID:    true,
			LogHost:         true,
			HandleError:     true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				fields := []zapcore.Field{
					zap.String("remote_ip", v.RemoteIP),
					zap.Duration("latency", v.Latency),
					zap.String("host", v.Host),
					zap.String("request", fmt.Sprintf("%s %s", v.Method, v.URI)),
					zap.Int("status", v.Status),
					zap.Int64("size", v.ResponseSize),
					zap.String("user_agent", v.UserAgent),
					zap.String("request_id", v.RequestID),
				}

				n := v.Status
				switch {
				case n >= 500:
					logger.With(zap.Error(v.Error)).Error("Server error", fields...)
				case n >= 400:
					logger.Info("Client error", fields...)
				case n >= 300:
					logger.Info("Redirection", fields...)
				default:
					logger.Info("Success", fields...)
				}

				return nil
			},
		}),
		sentryecho.New(sentryecho.Options{}),
	)

	return e
}
