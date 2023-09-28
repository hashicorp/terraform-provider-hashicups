package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
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
	Regions types.List         `tfsdk:"regions"`
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
			"regions": schema.ListAttribute{
				Description: "An Array of regions to look for subnets in.",
				ElementType: types.StringType,
				Optional:    true,
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

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Load AWS session configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to load AWS SDK config, "+ err.Error(), "")
		return
	}

	// Create an EC2 client
	client := ec2.NewFromConfig(cfg)

	var regions []string

	if state.Regions.IsNull() {
		// Get a list of all AWS regions
		describeRegionsResp, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})

		if err != nil {
			resp.Diagnostics.AddError("Error describing regions, " + err.Error(), "")
			return
		}
		for _, region := range describeRegionsResp.Regions {
		    regions = append(regions, *region.RegionName)
		}
	} else {
		resp.Diagnostics.Append(state.Regions.ElementsAs(ctx, &regions, false)...)
	}

	// Initialize variables for pagination
	var nextToken *string
	var subnetARNs []string

	// Iterate through regions
	for _, region := range regions {
		// Create a client for the current region
		regionCfg := cfg.Copy()
		regionCfg.Region = region
		regionClient := ec2.NewFromConfig(regionCfg)

		// Iterate through pages
		for {
			desribeSubnetsResp, err := regionClient.DescribeSubnets(
				ctx,
				&ec2.DescribeSubnetsInput{
					 Filters: []ec2types.Filter{
				        {
				            Name:   aws.String("map-public-ip-on-launch"),
				            Values: []string{"true"},
				        },
					},
					NextToken: nextToken,
				},
			)
			if err != nil {
				resp.Diagnostics.AddError("Error describing subnets in region, " + err.Error(), "")
				return
			}
			for _, subnet := range desribeSubnetsResp.Subnets {
				subnetARNs = append(subnetARNs, *subnet.SubnetArn)
			}

			// Check if there are more pages of subnets to retrieve
			if desribeSubnetsResp.NextToken == nil {
				break
			}

			nextToken = desribeSubnetsResp.NextToken
		}
	}
	// Sort lists
	sort.Strings(subnetARNs)
	sort.Strings(regions)

	// Ignore diags since we would catch error elsewhere
	state.Arns, _ = types.ListValueFrom(ctx, types.StringType, subnetARNs)
	state.Regions, _ = types.ListValueFrom(ctx, types.StringType, regions)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
