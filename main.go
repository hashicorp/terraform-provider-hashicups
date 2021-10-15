package main

import (
	"context"
	"terraform-provider-hashicups-pf/hashicups"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func main() {
	tfsdk.Serve(context.Background(), hashicups.New, tfsdk.ServeOpts{
		Name: "hashicups",
	})
}
