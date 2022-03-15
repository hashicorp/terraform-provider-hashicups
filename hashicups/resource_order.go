package hashicups

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.UseStateForUnknown(),
				},
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
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Order
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var items []hashicups.OrderItem
	for _, item := range plan.Items {
		items = append(items, hashicups.OrderItem{
			Coffee: hashicups.Coffee{
				ID: item.Coffee.ID,
			},
			Quantity: item.Quantity,
		})
	}

	// Create new order
	order, err := r.p.client.CreateOrder(items)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	// for more information on logging from providers, refer to
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	tflog.Trace(ctx, "created order", "order_id", order.ID)

	// Map response body to resource schema attribute
	var ois []OrderItem
	for _, oi := range order.Items {
		ois = append(ois, OrderItem{
			Coffee: Coffee{
				ID:          oi.Coffee.ID,
				Name:        types.String{Value: oi.Coffee.Name},
				Teaser:      types.String{Value: oi.Coffee.Teaser},
				Description: types.String{Value: oi.Coffee.Description},
				Price:       types.Number{Value: big.NewFloat(oi.Coffee.Price)},
				Image:       types.String{Value: oi.Coffee.Image},
			},
			Quantity: oi.Quantity,
		})
	}

	// Generate resource state struct
	var result = Order{
		ID:          types.String{Value: strconv.Itoa(order.ID)},
		Items:       ois,
		LastUpdated: types.String{Value: string(time.Now().Format(time.RFC850))},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceOrder) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Get current state
	var state Order
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get order from API and then update what is in state from what the API returns
	orderID := state.ID.Value

	// Get order current value
	order, err := r.p.client.GetOrder(orderID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading order",
			"Could not read orderID "+orderID+": "+err.Error(),
		)
		return
	}

	// Map response body to resource schema attribute
	state.Items = []OrderItem{}
	for _, item := range order.Items {
		state.Items = append(state.Items, OrderItem{
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
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceOrder) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	var plan Order
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state Order
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var items []hashicups.OrderItem
	for _, item := range plan.Items {
		items = append(items, hashicups.OrderItem{
			Coffee: hashicups.Coffee{
				ID: item.Coffee.ID,
			},
			Quantity: item.Quantity,
		})
	}

	// Get order ID from state
	orderID := state.ID.Value

	// Update order by calling API
	order, err := r.p.client.UpdateOrder(orderID, items)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error update order",
			"Could not update orderID "+orderID+": "+err.Error(),
		)
		return
	}

	// Map response body to resource schema attribute
	var ois []OrderItem
	for _, oi := range order.Items {
		ois = append(ois, OrderItem{
			Coffee: Coffee{
				ID:          oi.Coffee.ID,
				Name:        types.String{Value: oi.Coffee.Name},
				Teaser:      types.String{Value: oi.Coffee.Teaser},
				Description: types.String{Value: oi.Coffee.Description},
				Price:       types.Number{Value: big.NewFloat(oi.Coffee.Price)},
				Image:       types.String{Value: oi.Coffee.Image},
			},
			Quantity: oi.Quantity,
		})
	}

	// Generate resource state struct
	var result = Order{
		ID:          types.String{Value: strconv.Itoa(order.ID)},
		Items:       ois,
		LastUpdated: types.String{Value: string(time.Now().Format(time.RFC850))},
	}

	// Set state
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceOrder) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Order
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get order ID from state
	orderID := state.ID.Value

	// Delete order by calling API
	err := r.p.client.DeleteOrder(orderID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting order",
			"Could not delete orderID "+orderID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceOrder) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
