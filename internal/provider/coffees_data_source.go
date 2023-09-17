package provider

import (
	"context"
	"fmt"

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
	// Coffee coffeeModel `tfsdk:"coffees"`
	// AgentsIpv4               types.List `tfsdk:"agents_ipv4"`
	Arns types.List            `tfsdk:"arns"`
}

// types.List
// coffeesModel maps coffees schema data.
// type coffeeModel struct {
// 	Arns []types.String            `tfsdk:"arns"`
// }

// coffeesIngredientsModel maps coffee ingredients data.
type coffeesIngredientsModel struct {
	ID types.Int64 `tfsdk:"id"`
}

// Metadata returns the data source type name.
func (d *coffeesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_coffees"
}

// Schema defines the schema for the data source.
func (d *coffeesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve public subnet ARNs.",
		Attributes: map[string]schema.Attribute{
			"arns": schema.ListAttribute{
				Description: "An Array of public subnet ARNs.",
				ElementType: types.StringType,
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *coffeesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hashicups.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *coffeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state coffeesDataSourceModel

	arns := []string{"This", "is", "another", "list"}

	// TODO: Should we ignore diags here or not?
	// I think we should, as a response from Aws will contain more info
	state.Arns, _ = types.ListValueFrom(ctx, types.StringType, arns)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
