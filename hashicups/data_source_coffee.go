package hashicups

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceCoffeesType struct{}

func (r dataSourceCoffeesType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (r dataSourceCoffeesType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceCoffees{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceCoffees struct {
	p provider
}

func (r dataSourceCoffees) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	// Declare struct that this function will set to this data source's state
	var resourceState struct {
		Coffees []CoffeeIngredients `tfsdk:"coffees"`
	}

	coffees, err := r.p.client.GetCoffees()
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

	// Sample debug message
	// To view this message, set the TF_LOG environment variable to DEBUG
	// 		`export TF_LOG=DEBUG`
	// To hide debug message, unset the environment variable
	// 		`unset TF_LOG`
	fmt.Fprintf(stderr, "[DEBUG]-Resource State:%+v", resourceState)

	// Set state
	diags := resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
