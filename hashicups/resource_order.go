package hashicups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceOrderType struct{}

// Order Resource schema
func (r resourceOrderType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type: types.StringType,
				// When Computed is true, the provider will set value --
				// the user cannot define the value
				Computed: true,
			},
			"last_updated": {
				Type:     types.StringType,
				Computed: true,
			},
			"items": {
				// If Required is true, Terraform will throw error if user
				// doesn't specify value
				// If Optional is true, user can choose to supply a value
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"quantity": {
						Type:     types.NumberType,
						Required: true,
					},
					"coffee": {
						Required: true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"id": {
								Type:     types.NumberType,
								Required: true,
							},
							"name": {
								Type:     types.StringType,
								Computed: true,
							},
							"teaser": {
								Type:     types.StringType,
								Computed: true,
							},
							"description": {
								Type:     types.StringType,
								Computed: true,
							},
							"price": {
								Type:     types.NumberType,
								Computed: true,
							},
							"image": {
								Type:     types.StringType,
								Computed: true,
							},
						}),
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

// New resource instance
func (r resourceOrderType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceOrder{
		p: *(p.(*provider)),
	}, nil
}

type resourceOrder struct {
	p provider
}

// Create a new resource
func (r resourceOrder) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
}

// Read resource information
func (r resourceOrder) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {

}

// Update resource
func (r resourceOrder) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

// Delete resource
func (r resourceOrder) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
}

func (r resourceOrder) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
}
