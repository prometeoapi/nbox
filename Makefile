include ./Common.mk

GOBUILD=GO111MODULE=on CGO_ENABLED=0 installsuffix=cgo go build -trimpath

# ALL_MODULES includes ./* dirs (excludes . dir)
ALL_MODULES := $(shell find . -type f -name "go.mod" -exec dirname {} \; | sort | egrep  '^./' )

GOMODULES = $(ALL_MODULES) $(PWD)
TOOLS_MOD_DIR := $(abspath ./tools/workflow/linters)
TOOLS_BIN_DIR := $(abspath ./bin)

.PHONY: $(GOMODULES)
$(GOMODULES):
	@echo "Running target '$(TARGET)' in module '$@'"
	TOOL_BIN=$(TOOLS_BIN_DIR) $(MAKE) -C $@ $(TARGET)

# Triggers each module's delegation target
.PHONY: for-all-target
for-all-target: $(GOMODULES)

.PHONY: gofmt
gofmt:
	@$(MAKE) for-all-target TARGET="fmt"

.PHONY: golint
golint: lint-static-check
	@$(MAKE) for-all-target TARGET="lint"

.PHONY: gomod-tidy
gomod-tidy:
	@$(MAKE) for-all-target TARGET="mod-tidy"

.PHONY: gomod-vendor
gomod-vendor:
	go mod vendor

.PHONY: install-tools
install-tools:
	cd $(TOOLS_MOD_DIR) && GOBIN=$(TOOLS_BIN_DIR) go install golang.org/x/tools/cmd/goimports
	cd $(TOOLS_MOD_DIR) && GOBIN=$(TOOLS_BIN_DIR) go install honnef.co/go/tools/cmd/staticcheck
	cd $(TOOLS_MOD_DIR) && GOBIN=$(TOOLS_BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint
	cd $(TOOLS_MOD_DIR) && GOBIN=$(TOOLS_BIN_DIR) go install mvdan.cc/sh/v3/cmd/shfmt

.PHONY: install-all-deps
install-all-deps:
	go get ./...

.PHONY: lint-static-check
lint-static-check:
	@STATIC_CHECK_OUT=`$(TOOLS_BIN_DIR)/staticcheck $(ALL_PKGS) 2>&1`; \
		if [ "$$STATIC_CHECK_OUT" ]; then \
			echo "$(STATIC_CHECK) FAILED => static check errors:\n"; \
			echo "\033[0;31m$$STATIC_CHECK_OUT\033[0m\n"; \
			exit 1; \
		else \
			echo "\t- \033[0;32mStatic check finished successfully\033[0m"; \
		fi


.PHONY: imports-check
imports-check:
	@WARNINGS_CHECK_OUT=`$(TOOLS_BIN_DIR)/goimports -l .`; \
		if [ "$$WARNINGS_CHECK_OUT" ]; then \
			echo "Aborting commit due to bad code formatting.\n"; \
			echo "\033[0;31m$$WARNINGS_CHECK_OUT\033[0m\n"; \
			echo "You can use \"./bin/goimports -w .\" to automatically fix these issues or run."; \
			echo "\n\t make gofmt\n" ;\
			exit 1; \
		else \
			echo "\t- \033[0;32mGo imports check finished successfully\033[0m"; \
		fi

.PHONY: build
build:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ./build/darwin/amd64/microservice ./cmd/nbox
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ./build/linux/amd64/microservice ./cmd/nbox
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o ./build/linux/arm64/microservice ./cmd/nbox
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ./build/windows/amd64/microservice ./cmd/nbox

.PHONY: amd64-build
amd64-build:
#	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ./build/linux/amd64/microservice ./cmd/nbox
	GOOS=linux GOARCH=amd64 $(GOBUILD)  -o ./build/linux/amd64/microservice ./cmd/nbox


.PHONY: arm64-build
arm64-build:
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o ./build/linux/arm64/microservice ./cmd/nbox

.PHONY: clean
clean:
	rm -rf ./build


# Testing
.PHONY: test
test:
	go test ./... --cover

test-sonar:
	go test -covermode=atomic -coverprofile=coverage.out ./...
	go test -json ./... > report.json


run-local-sonar:
	docker run -d --name sonarqube -p 9000:9000 sonarqube