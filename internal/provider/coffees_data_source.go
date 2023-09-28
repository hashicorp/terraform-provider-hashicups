package provider

import (
	"context"
	"fmt"
	// "strings"
	"unsafe"

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
	Regions types.Set         `tfsdk:"regions"`
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
			"regions": schema.SetAttribute{
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
		resp.Diagnostics.AddError("WTF", "")
		return
	}

	// Load AWS session configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		resp.Diagnostics.AddError("Unable to load AWS SDK config, "+ err.Error(), "")
		return
	}

	// Create an EC2 client
	client := ec2.NewFromConfig(cfg)

	var regions []string
	// var shit *string

	// Trying out
	// hostFilters := make([]string, 0)
	var hostFilters []string
	// diags := diag.Diagnostics{}

	if state.Regions.IsNull() {
	// if true {
		// Get a list of all AWS regions
		describeRegionsResp, err := client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})

		if err != nil {
			resp.Diagnostics.AddError("Error describing regions, " + err.Error(), "")
			return
		}
		for _, region := range describeRegionsResp.Regions {
		    regions = append(regions, *region.RegionName)
		}
	} else {
		// regions = []string{"us-east-1", "us-east-2"}
		// state.Regions.ElementsAs(ctx, hostFilters, false)
		state.Regions.ElementsAs(ctx, &hostFilters, false)
		// diags.Append(state.Regions.ElementsAs(ctx, &hostFilters, false)...)
		// resp.Diagnostics.AddError(diags)
		// return

		// for _, element := range state.Regions.Elements() {
		// 	val, _ := element.ToTerraformValue(ctx)
		// 	foo := "hey"
		// 	shit = &foo
		// 	// resp.Diagnostics.AddError("Shit before is, "+ *shit, "")
		// 	err = val.As(shit)
		// 	if err != nil {
		// 		resp.Diagnostics.AddError("Hmm, "+ err.Error(), "")
		// 		return
		// 	}
		// 	// resp.Diagnostics.AddError("Shit after is, "+ *shit, "")
		// 	// return

		// 	// v, _ := val.value.(string)

		//     // regions = append(regions, value.String())
		//     // regions = append(regions, deepCopy(*shit))
		//     regions = append(regions, *shit)

		// 	// err := value.As(&regions)
		// 	// if err != nil {
		// 	// 	panic(err)
		// 	// }
		//     // regions = append(regions, )
		// }
	}

	// resp.Diagnostics.AddError("Shit hostFilters is, "+ hostFilters, "")
	// resp.Diagnostics.AddError("hostFilters is, "+ strings.Join(hostFilters, ","), "")
	// return

	regions = hostFilters

	// Initialize variables for pagination
	var nextToken *string
	var subnetARNs []string
	// regions = []string{"us-east-1", "us-east-2"}
	// regions = []string{"us-west-1", "us-east-1", "us-east-2"}

	// Iterate through regions
	// heepo := []string{"us-east-1", "us-east-2"}
	// for _, region := range heepo {
	for _, region := range regions {
		// Create a client for the current region
		regionCfg := cfg.Copy()
		regionCfg.Region = region
		regionClient := ec2.NewFromConfig(regionCfg)

		// Iterate through pages
		for {
			desribeSubnetsResp, err := regionClient.DescribeSubnets(
				context.TODO(),
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
				// TODO: WHAT DO THESE 2 ARGS DO?
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

	// TODO: Should we ignore diags here or not?
	// I think we should, as a response from Aws will contain more info if the Arns are messed up
	state.Arns, _ = types.ListValueFrom(ctx, types.StringType, subnetARNs)
	state.Regions, _ = types.SetValueFrom(ctx, types.StringType, regions)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}


func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}


func deepCopy(s string) string {
    b := make([]byte, len(s))
    copy(b, s)
    return *(*string)(unsafe.Pointer(&b))
}
