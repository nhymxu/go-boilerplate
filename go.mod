module github.com/nhymxu/go-boilerplate

go 1.26.2

replace rootPrj => ./

require (
	github.com/getsentry/sentry-go v0.46.1
	github.com/getsentry/sentry-go/echo v0.46.1
	github.com/knadh/koanf/parsers/dotenv v1.1.1
	github.com/knadh/koanf/providers/confmap v1.0.0
	github.com/knadh/koanf/providers/env v1.1.0
	github.com/knadh/koanf/providers/file v1.2.1
	github.com/knadh/koanf/v2 v2.3.4
	github.com/labstack/echo/v5 v5.1.0
	github.com/samber/slog-multi v1.8.0
	github.com/samber/slog-sentry/v2 v2.10.3
	github.com/spf13/cobra v1.10.2
	go.uber.org/automaxprocs v1.6.0
	rootPrj v0.0.0-00010101000000-000000000000
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/samber/lo v1.53.0 // indirect
	github.com/samber/slog-common v0.21.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	golang.org/x/mod v0.34.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/telemetry v0.0.0-20260311193753-579e4da9a98c // indirect
	golang.org/x/text v0.36.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	golang.org/x/tools v0.43.0 // indirect
)

tool golang.org/x/tools/cmd/deadcode
