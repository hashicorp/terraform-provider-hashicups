package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-provider-hashicups-pf/hashicups"
)

func main() {
	tfsdk.Serve(context.Background(), hashicups.New, tfsdk.ServeOpts{
		Name: "hashicups",
	})
}
