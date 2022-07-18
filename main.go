package main

import (
	"context"
	"terraform-provider-hashicups-pf/hashicups"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), hashicups.New, providerserver.ServeOpts{
		// NOTE: This is not a normal provider address, but it is used in
		// the example configurations and tutorial.
		Address: "hashicorp.com/edu/hashicups-pf",
	})
}
