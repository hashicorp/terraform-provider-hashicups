# Terraform Provider Scaffolding (Terraform Plugin Framework)

_This template repository expands on top of the https://github.com/hashicorp/terraform-provider-hashicups repo and provides a baseline for developing custom terraform providers_

This repository is a *template* for a [Terraform](https://www.terraform.io) provider. It is intended as a starting point for creating Terraform providers, containing:

- A resource and a data source (`internal/provider/`),
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.

These files contain boilerplate code that you will need to edit to create your own Terraform provider. Tutorials for creating Terraform providers can be found on the [HashiCorp Developer](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework) platform. _Terraform Plugin Framework specific guides are titled accordingly._

Please see the [GitHub template repository documentation](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) for how to create a new repository from this template on GitHub.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21
- [Docker](https://docs.docker.com/get-started/get-docker/)
- [AWS CLI](https://aws.amazon.com/cli/)
- [Make](https://www.gnu.org/software/make/)
- [mkcert](https://github.com/filosottile/mkcert)
- [gnupg](https://www.gnupg.org/)
- [jq](https://jqlang.org/)
- [goreleaser](https://goreleaser.com/)
- [golangci-lint](https://github.com/golangci/golangci-lint)


## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the provided Make command `build-and-package-local`:

```shell
make build-and-package-local
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Running a local Terraform Registry
The template provides setup configuration for a private Terraform registry that you can push and pull your custom providers to with [boring-registry](https://boring-registry.github.io/boring-registry/latest/) via docker compose.

Working with boring registry is largely abstracted via the Makefile, but a few one-time setup steps are required before first use.

### 1. Create a GPG signing key

Provider releases must be signed with a GPG key. If you don't already have one:

```shell
gpg --batch --gen-key <<EOF
%no-protection
Key-Type: RSA
Key-Length: 4096
Subkey-Type: RSA
Subkey-Length: 4096
Name-Real: Your Name
Name-Email: you@example.com
Expire-Date: 0
EOF
```

If your key has a passphrase, export it alongside the fingerprint (see the next step).

### 2. Configure `.env`

Create a `.env` file in the repository root. The Makefile picks this up automatically:

```shell
# Get your key fingerprint
gpg --list-secret-keys --with-colons --fingerprint | awk -F: '/^fpr:/ {print $10; exit}'
```

Then create `.env`:

```dotenv
GPG_FINGERPRINT=<your 40-character fingerprint>
GPG_PASSPHRASE=<your passphrase, or leave empty if the key has none>
```

### 3. Configure Terraform CLI credentials

Terraform needs a bearer token to authenticate against the local registry. Add the following to `~/.terraformrc`:

```hcl
credentials "localhost:5601" {
  token = "very-secure-token"
}
```

The token value must match `BORING_REGISTRY_AUTH_STATIC_TOKEN` in `docker/registry.docker-compose.yml`.

### 4. Add `minio` to `/etc/hosts`

The registry issues signed S3 download URLs using the internal Docker hostname `minio`. Terraform running on the host needs to resolve this name:

```shell
echo "127.0.0.1 minio" | sudo tee -a /etc/hosts
```

### Day-to-day workflow

Start the registry infrastructure (MinIO + boring-registry + nginx):

```shell
make start-local-registry-services
```

Build, sign, and publish the provider to the local registry:

```shell
make publish-to-local-registry
```

This single command:
1. Generates provider documentation
2. Builds binaries for all configured platforms
3. Signs the SHA256SUMS file with your GPG key
4. Generates `signing-keys.json` containing your GPG public key
5. Uploads all artifacts to MinIO under the correct boring-registry storage layout

To consume the published provider from a Terraform configuration, set the source to `localhost:5601/<namespace>/<provider>`:

```hcl
terraform {
  required_providers {
    hashicups = {
      source  = "localhost:5601/hashicups/hashicups"
      version = "0.0.1"
    }
  }
}
```

Then run `terraform init` as normal.

Stop the registry infrastructure when done:

```shell
make stop-local-registry-services
```

### Using published providers

After publishing a provider to the local registry, any Terraform configuration can consume it. This section guides both provider developers and end users on how to install and use published providers from the local boring-registry.

#### Prerequisites

Before using a published provider from the local registry, ensure:

1. **Registry infrastructure is running:**
   ```shell
   make start-local-registry-services
   ```
   This starts MinIO, boring-registry, and nginx. The registry is accessible at `https://localhost:5601`.

2. **Terraform CLI credentials are configured** in `~/.terraformrc`:
   ```hcl
   credentials "localhost:5601" {
     token = "very-secure-token"
   }
   ```
   The token must match `BORING_REGISTRY_AUTH_STATIC_TOKEN` in `docker/registry.docker-compose.yml` (currently `very-secure-token`).

3. **DNS resolution for MinIO is configured** in `/etc/hosts`:
   ```
   127.0.0.1 minio
   ```
   This allows Terraform to resolve signed S3 download URLs issued by boring-registry.

#### Configuring the provider in Terraform

In any Terraform configuration that needs the custom provider, add a `required_providers` block to the `terraform` block:

```hcl
terraform {
  required_providers {
    hashicups = {
      source  = "localhost:5601/hashicups/hashicups"
      version = "0.0.1"
    }
  }
}

provider "hashicups" {
  # Configure provider options here
}
```

**Key elements:**
- **source**: The provider source in the format `{hostname}/{namespace}/{provider-name}`
  - `localhost:5601` — Terraform registry hostname and port
  - `hashicups` — Namespace (organization)
  - `hashicups` — Provider name (customize for your provider)
- **version**: Must match a published version. Use `"0.0.1"` for the example, or `">= 0.0.1"` to accept newer versions

#### Complete example

Here's a minimal working Terraform configuration using the locally-published provider:

```hcl
terraform {
  required_providers {
    hashicups = {
      source  = "localhost:5601/hashicups/hashicups"
      version = "0.0.1"
    }
  }
}

provider "hashicups" {
  host     = "http://localhost:19090"  # HashiCups API endpoint
  username = "education"
  password = "test123"
}

# Use resources and data sources provided by the provider
data "hashicups_coffees" "example" {}

resource "hashicups_order" "example" {
  coffee {
    id       = data.hashicups_coffees.example.coffees[0].id
    quantity = 1
  }
}

output "order" {
  value = hashicups_order.example
}
```

#### Running Terraform with the local provider

1. **Initialize the working directory:**
   ```shell
   terraform init
   ```
   Terraform will authenticate using the credential in `~/.terraformrc`, download the provider from boring-registry, verify the GPG signature, and install it to `~/.terraform/plugins/`.

   Example output:
   ```
   Initializing the backend...
   Initializing provider plugins...
   - Finding latest version of localhost:5601/hashicups/hashicups...
   - Installing localhost:5601/hashicups/hashicups v0.0.1...
   - Installed localhost:5601/hashicups/hashicups v0.0.1 (self-signed, key ID 2035590153699EBE)
   
   Terraform has been successfully initialized!
   ```

2. **Use provider resources and data sources** normally:
   ```shell
   terraform plan
   terraform apply
   ```

#### Authentication troubleshooting

**Problem:** `Error: Failed to retrieve provider from registry: 401 Unauthorized`
- **Solution:** Verify the `credentials "localhost:5601"` block in `~/.terraformrc` and confirm the token matches `BORING_REGISTRY_AUTH_STATIC_TOKEN` in the docker-compose file.

**Problem:** `Error: provider binary not found`
- **Solution:** Ensure the provider was successfully published with `make publish-to-local-registry` and that the version in `required_providers` matches a published version. Check S3 contents:
  ```shell
  AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin \
    aws s3 ls s3://pvt-registry/providers/ --recursive --endpoint-url http://localhost:9000
  ```

**Problem:** `dial tcp: lookup minio: no such host`
- **Solution:** Add `127.0.0.1 minio` to `/etc/hosts`. The registry issues download URLs using the internal Docker hostname; the host machine needs to resolve this name.


