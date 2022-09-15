package main

import (
	"context"
	"terraform-provider-hashicups-pf/hashicups"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), hashicups.New, providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/hashicorp/hashicups. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "hashicorp.com/edu/hashicups-pf",
	})
}
