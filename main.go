package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-hashicups-pf/hashicups"
)

func main() {
	providerserver.Serve(context.Background(), hashicups.New, providerserver.ServeOpts{
		// NOTE: This is not a normal provider address, but it is used in
		// the example configurations and tutorial.
		Address: "hashicorp.com/edu/hashicups-pf",
	})
}
