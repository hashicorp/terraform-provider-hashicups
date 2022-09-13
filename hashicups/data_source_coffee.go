package hashicups

import (
	"context"
	"math/big"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satifies the expected interfaces.
var (
	_ datasource.DataSource              = &coffeesDataSource{}
	_ datasource.DataSourceWithConfigure = &coffeesDataSource{}
)

func NewCoffeesDataSource() datasource.DataSource {
	return &coffeesDataSource{}
}

type coffeesDataSource struct {
	client *hashicups.Client
}

func (d *coffeesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_coffees"
}

func (d *coffeesDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"coffees": {
				// When Computed is true, the provider will set value --
				// the user cannot define the value
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.NumberType,
						Computed: true,
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
					"ingredients": {
						Computed: true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"id": {
								Type:     types.NumberType,
								Computed: true,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (d *coffeesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*hashicups.Client)
}

func (d *coffeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Declare struct that this function will set to this data source's state
	var resourceState struct {
		Coffees []CoffeeIngredients `tfsdk:"coffees"`
	}

	coffees, err := d.client.GetCoffees()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving coffee",
			err.Error(),
		)
		return
	}

	// Map response body to resource schema
	for _, coffee := range coffees {
		c := CoffeeIngredients{
			ID:          coffee.ID,
			Name:        types.String{Value: coffee.Name},
			Teaser:      types.String{Value: coffee.Teaser},
			Description: types.String{Value: coffee.Description},
			Price:       types.Number{Value: big.NewFloat(coffee.Price)},
			Image:       types.String{Value: coffee.Image},
		}

		var ingredients []IngredientID
		for _, ingredient := range coffee.Ingredient {
			ingredients = append(ingredients, IngredientID{
				ID: ingredient.ID,
			})
		}

		c.Ingredient = ingredients

		resourceState.Coffees = append(resourceState.Coffees, c)
	}

	// for more information on logging from providers, refer to
	// https://terraform.io/plugin/log
	tflog.Trace(ctx, "Found coffees", map[string]any{"coffee_count": len(resourceState.Coffees)})

	// Set state
	diags := resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
