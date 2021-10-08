package hashicups

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Order -
type Order struct {
	ID          types.String `tfsdk:"id"`
	Items       []OrderItem  `tfsdk:"items"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// OrderItem -
type OrderItem struct {
	Coffee   Coffee `tfsdk:"coffee"`
	Quantity int    `tfsdk:"quantity"`
}

// Coffee -
// This Coffee struct is for Order.Items[].Coffee which does not have an
// ingredients field in the schema defined by plugin framework. Since the
// resource schema must match the struct exactly (extra field will return an
// error). This struct has Ingredients commented out.
type Coffee struct {
	ID          int          `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Teaser      types.String `tfsdk:"teaser"`
	Description types.String `tfsdk:"description"`
	Price       types.Number `tfsdk:"price"`
	Image       types.String `tfsdk:"image"`
	// Ingredients []Ingredient   `tfsdk:"ingredients"`
}

// Ingredient -
type Ingredient struct {
	ID       int    `tfsdk:"ingredient_id"`
	Name     string `tfsdk:"name"`
	Quantity int    `tfsdk:"quantity"`
	Unit     string `tfsdk:"unit"`
}

//
// Coffee Data Source specific structs
//

// CoffeeIngredients
type CoffeeIngredients struct {
	ID          int            `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Teaser      types.String   `tfsdk:"teaser"`
	Description types.String   `tfsdk:"description"`
	Price       types.Number   `tfsdk:"price"`
	Image       types.String   `tfsdk:"image"`
	Ingredient  []IngredientID `tfsdk:"ingredients"`
}

// Ingredient -
type IngredientID struct {
	ID int `tfsdk:"id"`
}
