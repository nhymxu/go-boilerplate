package api

import (
	"fmt"
	"net/http"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nhymxu/go-boilerplate/pkg/config"
)

func New() *echo.Echo {
	e := echo.New()
	if config.ENV.Debug {
		e.Debug = true
	}

	logger := zap.L()

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
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			Skipper: func(c echo.Context) bool {
				if c.Request().URL.Path == "/special-endpoint-can-replace-later" {
					return true
				}

				return false
			},
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
					logger.With(zap.Error(v.Error)).Warn("Client error", fields...)
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

	v1 := e.Group("/v1")
	groupV1Routes(v1)

	return e
}
