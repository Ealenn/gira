#COLORS
GREEN  := $(shell tput -Txterm setaf 2)
WHITE  := $(shell tput -Txterm setaf 7)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

# Set SHELL to bash for compatibility
SHELL := /bin/bash

# ---------------------------------------------------------------------
#                            HELP
# ---------------------------------------------------------------------
HELP_FUN = \
	%help; \
	while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^([a-zA-Z\-]+)\s*:.*\#\#(?:@([a-zA-Z\-]+))?\s(.*)$$/ }; \
	print "usage: make [target]\n\n"; \
	for (sort keys %help) { \
	print "${WHITE}$$_:${RESET}\n"; \
	for (@{$$help{$$_}}) { \
	$$sep = " " x (32 - length $$_->[0]); \
	print "  ${YELLOW}$$_->[0]${RESET}$$sep${GREEN}$$_->[1]${RESET}\n"; \
	}; \
	print "\n"; }

.PHONY: help
help: ##@other Display available commands and their descriptions
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

# ---------------------------------------------------------------------
#                            PROJECT
# ---------------------------------------------------------------------
.PHONY: build
build: tidy fmt ##@project Compile Gira binary
	@go build -ldflags="-s -w" ./cmd/...

.PHONY: test
test: ##@project Run all unit tests
	@go test ./... -v

.PHONY: run
run: ##@project Execute Gira with provided arguments
	@go run ./cmd/gira $(filter-out $@, $(MAKECMDGOALS))

.PHONY: debug
debug: ##@project Execute Gira in verbose mode with provided arguments
	@go run ./cmd/gira $(filter-out $@, $(MAKECMDGOALS)) --verbose

# ---------------------------------------------------------------------
#                            TOOLS
# ---------------------------------------------------------------------
.PHONY: tidy
tidy: ##@tools Remove unused Go modules and tidy dependencies
	@go mod tidy

.PHONY: upgrade
upgrade: tidy ##@tools Upgrade all dependencies to the latest version
	@go get -u ./...

.PHONY: lint
lint: ##@tools Run linter to check code style
	@golangci-lint run ./...

.PHONY: fmt
fmt: ##@tools Format Go source code
	@go fmt ./...

.PHONY: vet
vet: ##@tools Run Go vet to check for suspicious code
	@go vet ./...
