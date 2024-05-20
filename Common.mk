# commons --

# ALL_PKGS is the list of all packages where ALL_SRC files reside.
ALL_PKGS := $(sort $(shell go list ./...))
APP_IMPORT_PATH=proteccion/replica

GOTEST_OPT?= -race -timeout 120s
GOCMD?= go
GO_ACC=go-acc
LINT=golangci-lint

GOOS := $(shell $(GOCMD) env GOOS)
GOARCH := $(shell $(GOCMD) env GOARCH)

#GOIMPORTS_OPT?= -w -local .
GOIMPORTS_OPT?= -w .


.PHONY: fmt
fmt:
	go fmt ./...
	$(TOOL_BIN)/goimports $(GOIMPORTS_OPT) ./

.PHONY: lint
lint:
	$(TOOL_BIN)/$(LINT) run --timeout 5m --enable gosec

.PHONY: mod-tidy
mod-tidy:
	go mod tidy
