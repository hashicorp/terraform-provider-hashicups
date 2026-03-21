# Variables
VERSION := 0.0.1
PROJECT_NAME := hashicups
CERT_DIR := certs
CERT_FILE := $(CERT_DIR)/cert.pem
KEY_FILE := $(CERT_DIR)/key.pem

# Help target
GREY := $(shell tput setaf 8)
GREEN_BOLD := $(shell tput setaf 2; tput bold)
RESET := $(shell tput sgr0)
##@ Help
help: ## Show this help message
	@echo "Usage: make [target]"
	@awk 'BEGIN {FS = ":.*##";} \
	/^[a-zA-Z_-]+:.*##/ { printf "  ${GREY}%-40s${RESET} %s\n", $$1, $$2 } \
	/^##@/ { printf "\n${GREEN_BOLD}%-40s${RESET}\n", substr($$0, 5) }' \
	$(MAKEFILE_LIST)

##@ Docker Compose
start-local-hashicups: ## Start local HashiCups server
	@echo "Starting local HashiCups server ..."
	docker compose --file docker/hashicups.docker-compose.yml up --remove-orphans

stop-local-hashicups: ## Stop local HashiCups server
	@echo "Stopping local HashiCups server ..."
	docker compose --file docker/hashicups.docker-compose.yml down --remove-orphans

$(CERT_FILE):
	@echo "Certs missing. Generating self-signed certificate for testing..."
	mkdir -p $(CERT_DIR)
	mkcert -install
	mkcert -key-file $(KEY_FILE) -cert-file $(CERT_FILE) localhost 127.0.0.1
	@echo "Certificate generated at $(CERT_FILE) and key at $(KEY_FILE)"

start-local-registry-services: $(CERT_FILE) ## Start local terraform registry services to host binaries
	@echo "Starting local terraform registry services to host binaries..."
	docker compose --file docker/registry.docker-compose.yml up --remove-orphans

stop-local-registry-services: ## Stop local terraform registry services
	@echo "Stopping local terraform registry services..."
	docker compose --file docker/registry.docker-compose.yml down --remove-orphans

##@ Development
tidy: ## Clean up go.mod and go.sum
	go mod tidy

fmt: ## Format code and apply automatic fixes
	gofmt -s -w -e .
	golangci-lint run --fix

lint: fmt ## Run linters
	golangci-lint run

test: lint ## Run unit tests
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc: lint ## Run acceptance tests
	TF_ACC=1 go test -v -cover -timeout 120m ./...

##@ Build and Install
generate-docs: ## Generate documentation
	cd tools; go generate ./...

build-and-package-local: generate-docs ## Build and package the provider locally without signing
	PROJECT_NAME=$(PROJECT_NAME) \
	BUILD_VERSION=$(VERSION) \
	goreleaser release --snapshot --clean --skip=sign

build-and-package: generate-docs ## Build and package the provider
	PROJECT_NAME=$(PROJECT_NAME) \
	BUILD_VERSION=$(VERSION) \
	goreleaser release --snapshot --clean

##@ Clean up
clean: ## Clean up generated files
	rm -rf dist/ bin/ $(CERT_DIR)/

.PHONY: help start-local-hashicups stop-local-hashicups start-local-registry-services stop-local-registry-services tidy fmt lint test testacc generate-docs build-and-package clean
