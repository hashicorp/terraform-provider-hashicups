# Variables
CERT_DIR := certs
CERT_FILE := $(CERT_DIR)/cert.pem
KEY_FILE := $(CERT_DIR)/key.pem


# Targets
$(CERT_FILE):
	@echo "Certs missing. Generating self-signed certificate for testing..."
	mkdir -p $(CERT_DIR)
	mkcert -install
	mkcert -key-file $(KEY_FILE) -cert-file $(CERT_FILE) localhost 127.0.0.1
	@echo "Certificate generated at $(CERT_FILE) and key at $(KEY_FILE)"

start-local-terraform-registry: $(CERT_FILE)
	@echo "Starting local Terraform registry with self-signed certificate..."
	docker compose --file docker/registry.docker-compose.yml up --remove-orphans

stop-local-terraform-registry:
	@echo "Stopping local Terraform registry..."
	docker compose --file docker/registry.docker-compose.yml down --remove-orphans

default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
