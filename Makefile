.DEFAULT_GOAL := list

# Insert a comment starting with '##' after a target, and it will be printed by 'make' and 'make list'
.PHONY: list
list: ## list Makefile targets
	@echo "The most used targets: \n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: check-fmt
check-fmt: ## Ensure code is formatted
	gofmt -l -d . 	# For the sake of debugging
	test -z "$$(gofmt -l .)"

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: tests
tests: ## Run all tests and requires a running rabbitmq-server. Use GO_TEST_FLAGS to add extra flags to go test
	go test -race -v -tags integration $(GO_TEST_FLAGS)

.PHONY: check
check:
	#golangci-lint run ./...
	./scripts/lint.sh

## tidy: format code and tidy modfile
.PHONY: tidy
tidy: fmt
	go mod tidy -v

.PHONY: find-cgo-pkg
find-cgo-pkg:  ## identify which package on project using CGO
	./scripts/find-cgo-pkg.sh

.PHONE: check-duplicate-code
check-duplicate-code: ## identify duplicate code inside a project
	go install github.com/boyter/dcd@latest
	dcd
