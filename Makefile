include .env
include .sercet.env
export

.PHONY: help
help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# develop

.PNONY: env
env: ## Create .env file.
	@cp .env.local .env

.PHONY: aqua
aqua: ## Put the path in your environment variables. ex) export PATH="${AQUA_ROOT_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/aquaproj-aqua}/bin:$PATH"
	@go run github.com/aquaproj/aqua-installer@latest --aqua-version v2.0.0

.PHONY: tool
tool: ## Install tool.
	@aqua i
	@(cd api && npm install)

# go

.PHONY: fmt
fmt: ## Format code.
	@go fmt ./...

.PHONY: lint
lint: ## Lint code. FIXME: https://github.com/golangci/golangci-lint/issues/3711
	@golangci-lint run ./... --fix

.PHONY: mod
mod: ## Go mod
	@go mod tidy
	@go mod vendor

.PHONY: modules
modules: ## List modules with dependencies.
	@go list -u -deps ./...

.PHONY: renovate
renovate: ## Update modules with dependencies.
	@go get -u -t ./...
	@go mod tidy
	@go mod vendor

.PHONY: compile
compile: ## Compile code.
	@go build -v ./... && go clean

.PHONY: test
test: ## Run unit test. If you want to invalidate the cache, please specify an argument like `make test c=c`.
	@$(call _test,${c})

define _test
if [ -z "$1" ]; then \
	go test ./... ; \
else \
	go test ./... -count=1 ; \
fi
endef

.PHONY: buflint
buflint: ## Lint proto file.
	@(cd proto && buf lint)

.PHONY: bufmt
bufmt: ## Format proto file.
	@(cd proto && buf format -w)

.PHONY: apilint
apilint: ## Lint api file.
	@(cd api && npx spectral lint openapi.yaml)

.PHONY: ymlint
ymlint: ## Lint yaml file.
	@yamlfmt -lint

.PHONY: ymlfmt
ymlfmt: ## Format yaml file.
	@yamlfmt
