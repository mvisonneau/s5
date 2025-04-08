NAME          := s5
FILES         := $(shell git ls-files */*.go)
COVERAGE_FILE := coverage.out
REPOSITORY    := mvisonneau/$(NAME)

GOLANGCI_LINT := go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint
GORELEASER    := go tool github.com/goreleaser/goreleaser/v2

.DEFAULT_GOAL := help

.PHONY: fmt
fmt: ## Format source code
	$(GOLANGCI_LINT) fmt -v

.PHONY: lint
lint: ## Run all lint related tests upon the codebase
	$(GOLANGCI_LINT) run -v

.PHONY: test
test: ## Run the tests against the codebase
	@rm -rf $(COVERAGE_FILE)
	go test -v -count=1 -race ./... -coverprofile=$(COVERAGE_FILE)
	@go tool cover -func $(COVERAGE_FILE) | awk '/^total/ {print "coverage: " $$3}'

.PHONY: coverage
coverage: ## Prints coverage report
	go tool cover -func $(COVERAGE_FILE)

.PHONY: install
install: ## Build and install locally the binary (dev purpose)
	go install ./cmd/$(NAME)

.PHONY: build
build: ## Build the binaries using local GOOS
	go build ./cmd/$(NAME)

.PHONY: release
release: ## Build & release the binaries (stable)
	mkdir -p ${HOME}/.cache/snapcraft/download
	mkdir -p ${HOME}/.cache/snapcraft/stage-packages
	git tag -d edge
	$(GORELEASER) release --clean
	find dist -type f -name "*.snap" -exec snapcraft upload --release stable,edge '{}' \;

.PHONY: prerelease
prerelease: ## Build & prerelease the binaries (edge)
	@\
		REPOSITORY=$(REPOSITORY) \
		NAME=$(NAME) \
		GITHUB_TOKEN=$(GITHUB_TOKEN) \
		.github/prerelease.sh

.PHONY: clean
clean: ## Remove binary if it exists
	rm -f $(NAME)

.PHONY: coverage-html
coverage-html: ## Generates coverage report and displays it in the browser
	go tool cover -html=coverage.out

.PHONY: dev-env
dev-env: ## Build a local development environment using Docker
	@docker run -d --cap-add IPC_LOCK --name vault vault:$(VAULT_VERSION)
	@sleep 2
	@docker exec -it \
		-e VAULT_ADDR=http://localhost:8200 \
		-e VAULT_TOKEN=$$(docker logs vault 2>/dev/null | grep 'Root Token' | cut -d' ' -f3 | sed -E "s/[[:cntrl:]]\[[0-9]{1,3}m//g") \
		vault \
		/bin/sh -c "vault secrets enable transit"
	@docker run -it --rm \
		-v $(shell pwd):/go/src/github.com/mvisonneau/$(NAME) \
		-w /go/src/github.com/mvisonneau/$(NAME) \
		-e VAULT_ADDR=http://$$(docker inspect vault | jq -r '.[0].NetworkSettings.IPAddress'):8200 \
		-e VAULT_TOKEN=$$(docker logs vault 2>/dev/null | grep 'Root Token' | cut -d' ' -f3 | sed -E "s/[[:cntrl:]]\[[0-9]{1,3}m//g") \
		-e S5_TRANSIT_KEY=foo \
		goreleaser/goreleaser:v1.24.0 \
		/bin/bash -c 'apk add --no-cache make; make setup; make install; bash'
	@docker kill vault
	@docker rm vault -f

.PHONY: is-git-dirty
is-git-dirty: ## Tests if git is in a dirty state
	@git status --porcelain
	@test $(shell git status --porcelain | grep -c .) -eq 0

.PHONY: man-pages
man-pages: ## Generates man pages
	rm -rf helpers/manpages
	mkdir -p helpers/manpages
	go run ./cmd/tools/man | gzip -c -9 >helpers/manpages/$(NAME).1.gz

.PHONY: autocomplete-scripts
autocomplete-scripts: ## Download CLI autocompletion scripts
	rm -rf helpers/autocomplete
	mkdir -p helpers/autocomplete
	curl -sL https://raw.githubusercontent.com/urfave/cli/v2.27.1/autocomplete/bash_autocomplete > helpers/autocomplete/bash
	curl -sL https://raw.githubusercontent.com/urfave/cli/v2.27.1/autocomplete/zsh_autocomplete > helpers/autocomplete/zsh

.PHONY: all
all: lint test build coverage ## Test, builds and ship package for all supported platforms

.PHONY: help
help: ## Displays this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
