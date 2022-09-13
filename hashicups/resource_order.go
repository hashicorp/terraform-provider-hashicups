package hashicups

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &orderResource{}
	_ resource.ResourceWithConfigure   = &orderResource{}
	_ resource.ResourceWithImportState = &orderResource{}
)

func NewOrderResource() resource.Resource {
	return &orderResource{}
}

type orderResource struct {
	client *hashicups.Client
}

func (r *orderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_order"
}

// Order Resource schema
func (r *orderResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				}),
			},
		},
	}, nil
}

// Configure resource instance
func (r *orderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*hashicups.Client)
}

// Create a new resource
func (r *orderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not configured",
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
	order, err := r.client.CreateOrder(items)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	// for more information on logging from providers, refer to
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	tflog.Trace(ctx, "created order", map[string]interface{}{"order_id": order.ID})

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
func (r *orderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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
	order, err := r.client.GetOrder(orderID)
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
func (r *orderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete resource
func (r *orderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// Import resource
func (r *orderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Save the import identifier in the id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
