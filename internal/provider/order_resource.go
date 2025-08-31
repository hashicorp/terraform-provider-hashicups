package provider

import (
	"context"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &orderResource{}
	_ resource.ResourceWithConfigure = &orderResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewOrderResource() resource.Resource {
	return &orderResource{}
}

// orderResourceModel maps the resource schema data.
type orderResourceModel struct {
	ID          types.String     `tfsdk:"id"`
	Items       []orderItemModel `tfsdk:"items"`
	LastUpdated types.String     `tfsdk:"last_updated"`
}

// orderItemModel maps order item data.
type orderItemModel struct {
	Coffee   orderItemCoffeeModel `tfsdk:"coffee"`
	Quantity types.Int64          `tfsdk:"quantity"`
}

// orderItemCoffeeModel maps coffee order item data.
type orderItemCoffeeModel struct {
	ID          types.Int64   `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	Teaser      types.String  `tfsdk:"teaser"`
	Description types.String  `tfsdk:"description"`
	Price       types.Float64 `tfsdk:"price"`
	Image       types.String  `tfsdk:"image"`
}

// orderResource is the resource implementation.
type orderResource struct {
	client *hashicups.Client
}

// Metadata returns the resource type name.
func (r *orderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_order"
}

// Schema defines the schema for the resource.
func (r *orderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an order.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the order.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the order.",
				Computed:    true,
			},
			"items": schema.ListNestedAttribute{
				Description: "List of items in the order.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"quantity": schema.Int64Attribute{
							Description: "Count of this item in the order.",
							Required:    true,
						},
						"coffee": schema.SingleNestedAttribute{
							Description: "Coffee item in the order.",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Description: "Numeric identifier of the coffee.",
									Required:    true,
								},
								"name": schema.StringAttribute{
									Description: "Product name of the coffee.",
									Computed:    true,
								},
								"teaser": schema.StringAttribute{
									Description: "Fun tagline for the coffee.",
									Computed:    true,
								},
								"description": schema.StringAttribute{
									Description: "Product description of the coffee.",
									Computed:    true,
								},
								"price": schema.Float64Attribute{
									Description: "Suggested cost of the coffee.",
									Computed:    true,
								},
								"image": schema.StringAttribute{
									Description: "URI for an image of the coffee.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Create a new resource.
func (r *orderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read resource information.
func (r *orderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *orderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

func (r *orderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	return
}

// Configure adds the provider configured client to the resource.
func (r *orderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	return
}
