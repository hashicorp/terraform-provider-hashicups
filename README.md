# Terraform Provider Hashicups

This repo is a companion repo to the [Call APIs with Terraform Providers](https://learn.hashicorp.com/collections/terraform/providers) Learn collection. 

In the collection, you will use the HashiCups provider as a bridge between Terraform and the HashiCups API. Then, extend Terraform by recreating the HashiCups provider. By the end of this collection, you will be able to take these intuitions to create your own custom Terraform provider. 

## Build & test the provider

Run the following commands to build and test the provider

```shell
$ make build

$ make test

# and, for acceptance tests
$ make testacc
```

## Test `examples/` configuration

First, build and install the provider.

```shell
$ make install
```

NOTE: Currently the `Makefile` is setup to `install` the built prodiver under `darwin_amd64` os/arch.
You might need to tweak that if you are on a different system (ex. `darwin_arm64` for Apple M1 hardware).

Then, launch [HashiCups](https://github.com/hashicorp-demoapp) locally, by leveraging 
[Docker Compose](https://docs.docker.com/compose/):

```shell
$ cd docker_compose
$ docker-compose up -d
```

NOTE: When launching this setup for the first time, there is a one time operation required - creating 
creadentials in HashiCups. These are the cretentials that we will use when executing terraform against it.
To create them:

```shell
curl -X POST localhost:19090/signup -d '{"username":"education", "password":"test123"}'
```

Once HashiCups is running and the user `education` is created, it's time to launch the example:

```shell
$ cd examples
$ terraform init
$ terraform plan
...
$ terraform apply
```

NOTE: Every time a new build of the provider gets `make install`-ed, you will need to delete
the `.terraform.lock.hcl` file inside the `examples/` directory. If not, `terraform` will refuse to `init` 
because a checksum mismatch will be detected.