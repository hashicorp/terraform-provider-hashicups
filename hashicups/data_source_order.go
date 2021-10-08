package hashicups

import (
	"context"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceOrderType struct{}

func (r dataSourceOrderType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type: types.StringType,
				// When Computed is true, the provider will set value --
				// the user cannot define the value
				Required: true,
			},
			"last_updated": {
				Type:     types.StringType,
				Computed: true,
			},
			"items": {
				// If Required is true, Terraform will throw error if user
				// doesn't specify value
				// If Optional is true, user can choose to supply a value
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"quantity": {
						Type:     types.NumberType,
						Computed: true,
					},
					"coffee": {
						Computed: true,
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

func (r dataSourceOrderType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceOrder{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceOrder struct {
	p provider
}

func (r dataSourceOrder) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {

	// Declare struct that this function will set to this data source's state

	var resourceData Order
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get order from API and then update what is in state from what the API returns
	orderID := resourceData.ID.Value

	// Get order current value
	order, err := r.p.client.GetOrder(orderID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving order",
			err.Error(),
		)
		return
	}

	// Map response body to resource schema attribute
	resourceData.Items = []OrderItem{}
	for _, item := range order.Items {
		resourceData.Items = append(resourceData.Items, OrderItem{
			Coffee: Coffee{
				ID:          item.Coffee.ID,
				Name:        types.String{Value: item.Coffee.Name},
				Teaser:      types.String{Value: item.Coffee.Teaser},
				Description: types.String{Value: item.Coffee.Description},
				Price:       types.Number{Value: big.NewFloat(item.Coffee.Price)},
				Image:       types.String{Value: item.Coffee.Image},
			},
			Quantity: item.Quantity,
		})
	}

	// Set state
	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
