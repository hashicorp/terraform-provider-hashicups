package hashicups

import (
	"context"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &coffeesDataSource{}
	_ datasource.DataSourceWithConfigure = &coffeesDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewCoffeesDataSource() datasource.DataSource {
	return &coffeesDataSource{}
}

// coffeesDataSource is the data source implementation.
type coffeesDataSource struct {
	client *hashicups.Client
}

// coffeesDataSourceModel maps the data source schema data.
type coffeesDataSourceModel struct {
	Coffees []coffeesModel `tfsdk:"coffees"`
	ID      types.String   `tfsdk:"id"`
}

// coffeesModel maps coffees schema data.
type coffeesModel struct {
	ID          types.Int64               `tfsdk:"id"`
	Name        types.String              `tfsdk:"name"`
	Teaser      types.String              `tfsdk:"teaser"`
	Description types.String              `tfsdk:"description"`
	Price       types.Float64             `tfsdk:"price"`
	Image       types.String              `tfsdk:"image"`
	Ingredients []coffeesIngredientsModel `tfsdk:"ingredients"`
}

// coffeesIngredientsModel maps coffee ingredients data
type coffeesIngredientsModel struct {
	ID types.Int64 `tfsdk:"id"`
}

// Metadata returns the data source type name.
func (d *coffeesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_coffees"
}

// GetSchema defines the schema for the data source.
func (d *coffeesDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Fetches the list of coffees.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Placeholder identifier attribute.",
				Computed:    true,
				Type:        types.StringType,
			},
			"coffees": {
				Description: "List of coffees.",
				Computed:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Description: "Numeric identifier of the coffee.",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"name": {
						Description: "Product name of the coffee.",
						Type:        types.StringType,
						Computed:    true,
					},
					"teaser": {
						Description: "Fun tagline for the coffee.",
						Type:        types.StringType,
						Computed:    true,
					},
					"description": {
						Description: "Product description of the coffee.",
						Type:        types.StringType,
						Computed:    true,
					},
					"price": {
						Description: "Suggested cost of the coffee.",
						Type:        types.Float64Type,
						Computed:    true,
					},
					"image": {
						Description: "URI for an image of the coffee.",
						Type:        types.StringType,
						Computed:    true,
					},
					"ingredients": {
						Description: "List of ingredients in the coffee.",
						Computed:    true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"id": {
								Description: "Numeric identifier of the coffee ingredient.",
								Type:        types.Int64Type,
								Computed:    true,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

// Configure adds the provider configured client to the data source.
func (d *coffeesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*hashicups.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *coffeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state coffeesDataSourceModel

	coffees, err := d.client.GetCoffees()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Coffees",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, coffee := range coffees {
		coffeeState := coffeesModel{
			ID:          types.Int64Value(int64(coffee.ID)),
			Name:        types.StringValue(coffee.Name),
			Teaser:      types.StringValue(coffee.Teaser),
			Description: types.StringValue(coffee.Description),
			Price:       types.Float64Value(coffee.Price),
			Image:       types.StringValue(coffee.Image),
		}

		for _, ingredient := range coffee.Ingredient {
			coffeeState.Ingredients = append(coffeeState.Ingredients, coffeesIngredientsModel{
				ID: types.Int64Value(int64(ingredient.ID)),
			})
		}

		state.Coffees = append(state.Coffees, coffeeState)
	}

	state.ID = types.StringValue("placeholder")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
