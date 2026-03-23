# Load environment variables from .env file if it exists
-include .env
export

# Variables
VERSION := 0.0.1
PROJECT_NAME := hashicups
NAMESPACE := hashicups
PROVIDER_NAME := hashicups
S3_BUCKET := pvt-registry
S3_ENDPOINT := http://localhost:9000
GPG := gpg
DEFAULT_GPG_FINGERPRINT := $(shell $(GPG) --list-secret-keys --with-colons --fingerprint 2>/dev/null | awk -F: '/^fpr:/ {print $$10; exit}')
CERT_DIR := certs
CERT_FILE := $(CERT_DIR)/cert.pem
KEY_FILE := $(CERT_DIR)/key.pem
AWS_ACCESS_KEY_ID := minioadmin
AWS_SECRET_ACCESS_KEY := minioadminpassword
AWS_DEFAULT_REGION := eu-east-1

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

publish-to-local-registry: build-and-package generate-signing-keys upload-to-registry ## Publish provider binaries to local terraform registry. Requires GPG signing!
	@echo "Publishing provider binaries to local terraform registry..."
	@echo "✓ Provider published to MinIO S3 backend"

generate-signing-keys: check-gpg-signing ## Generate signing-keys.json for the terraform registry
	@mkdir -p dist; \
	GPG_KEY=$${GPG_FINGERPRINT:-$(DEFAULT_GPG_FINGERPRINT)}; \
	KEY_ID=$$($(GPG) --list-keys --with-colons $$GPG_KEY 2>/dev/null | awk -F: '/^pub:/ {print substr($$5, length($$5)-15); exit}'); \
	$(GPG) --armor --export $$GPG_KEY 2>/dev/null | \
	jq -Rs --arg key_id "$$KEY_ID" '{"gpg_public_keys": [{"key_id": $$key_id, "ascii_armor": .}]}' > dist/signing-keys.json; \
	echo "Generated dist/signing-keys.json with key ID: $$KEY_ID"

upload-to-registry: ## Upload generated provider files to MinIO S3 backend following boring-registry layout
	@echo "Uploading provider to S3 backend..."
	@command -v aws >/dev/null 2>&1 || (echo "Error: aws CLI not found. Install it from https://aws.amazon.com/cli/" && exit 1)
	@DIST_DIR=dist; \
	if [ ! -d "$$DIST_DIR" ]; then echo "Error: dist directory not found. Run 'make build-and-package' first."; exit 1; fi; \
	echo "Uploading signing-keys.json to providers/$(NAMESPACE)/"; \
	aws s3 cp $$DIST_DIR/signing-keys.json s3://$(S3_BUCKET)/providers/$(NAMESPACE)/signing-keys.json \
		--endpoint-url $(S3_ENDPOINT) \
		--region $(AWS_DEFAULT_REGION) || (echo "Error uploading signing keys"; exit 1); \
	echo "Uploading provider files to providers/$(NAMESPACE)/$(PROVIDER_NAME)/"; \
	for file in $$DIST_DIR/terraform-provider-$(PROVIDER_NAME)_$(VERSION)_*.zip $$DIST_DIR/terraform-provider-$(PROVIDER_NAME)_$(VERSION)_SHA256SUMS*; do \
		if [ -f "$$file" ]; then \
			echo "  Uploading $$(basename $$file)"; \
			aws s3 cp "$$file" "s3://$(S3_BUCKET)/providers/$(NAMESPACE)/$(PROVIDER_NAME)/$$(basename $$file)" \
				--endpoint-url $(S3_ENDPOINT) \
				--region $(AWS_DEFAULT_REGION) || (echo "Error uploading $$file"; exit 1); \
		fi; \
	done; \
	echo "Upload complete!"

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

check-gpg-signing: ## Verify GPG signing prerequisites
	@command -v $(GPG) >/dev/null 2>&1 || (echo "Error: gpg not found on PATH." && exit 1)
	@if [ -z "$(GPG_FINGERPRINT)" ] && [ -z "$(DEFAULT_GPG_FINGERPRINT)" ]; then \
		echo "Error: no GPG secret key found. Import/create one or set GPG_FINGERPRINT."; \
		exit 1; \
	fi
	@GPG_KEY=$${GPG_FINGERPRINT:-$(DEFAULT_GPG_FINGERPRINT)}; \
	echo "Using GPG fingerprint: $$GPG_KEY"; \
	TMP_FILE=$$(mktemp); \
	echo "gpg-sign-check" > "$$TMP_FILE"; \
	if ! $(GPG) --batch --yes --pinentry-mode loopback --passphrase "$${GPG_PASSPHRASE:-}" --local-user "$$GPG_KEY" --output /dev/null --detach-sign "$$TMP_FILE" >/dev/null 2>&1; then \
		rm -f "$$TMP_FILE"; \
		echo "Error: unable to sign with the selected GPG key. If the key is passphrase-protected, export GPG_PASSPHRASE."; \
		exit 1; \
	fi; \
	rm -f "$$TMP_FILE"

build-and-package: generate-docs check-gpg-signing ## Build and package the provider with GPG signing
	@GPG_FINGERPRINT=$${GPG_FINGERPRINT:-$(DEFAULT_GPG_FINGERPRINT)} \
	GPG_PASSPHRASE=$${GPG_PASSPHRASE:-} \
	PROJECT_NAME=$(PROJECT_NAME) \
	BUILD_VERSION=$(VERSION) \
	goreleaser release --snapshot --clean

##@ Clean up
clean: ## Clean up generated files
	rm -rf dist/ bin/ $(CERT_DIR)/

.PHONY: help start-local-hashicups stop-local-hashicups start-local-registry-services publish-to-local-registry stop-local-registry-services tidy fmt lint test testacc generate-docs generate-signing-keys upload-to-registry check-gpg-signing build-and-package-local build-and-package clean
