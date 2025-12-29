.PHONY: fmt git-hooks lint test down build dev-up up dev logs clean

GOOS ?= linux
GOARCH ?= amd64
GO ?= CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go

LINT_VERSION := v2.7.2
GOLANGCI_LINT := $$(go env GOPATH)/bin/golangci-lint

fmt:
	$(GO) fmt ./...

git-hooks:
	git config --local core.hooksPath .githooks/

lint: fmt
	@if ! [ -x "$(GOLANGCI_LINT)" ]; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
			| sh -s -- -b "$$(go env GOPATH)/bin" $(LINT_VERSION) ; \
	fi
	@echo "linting..."
	@"$(GOLANGCI_LINT)" run --timeout=5m
	@echo "linting completed!"

test:
	go test -v -race ./...

down:
	docker-compose down -v

build:
	docker-compose build

dev-up:
	docker-compose up

up:
	docker-compose up -d

dev: down build dev-up
	@echo "app is running"
	@echo "check logs via make logs"

logs:
	docker-compose logs -f

clean: down
	docker system prune -f
