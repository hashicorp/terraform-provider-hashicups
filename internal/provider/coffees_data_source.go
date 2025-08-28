package provider

import (
	"context"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

// Schema defines the schema for the data source.
func (d *coffeesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{}

	return
}

// Read refreshes the Terraform state with the latest data.
func (d *coffeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	return
}
