package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &publicELBsDataSource{}
	_ datasource.DataSourceWithConfigure = &publicELBsDataSource{}
)

// NewPublicEC2sDataSource is a helper function to simplify the provider implementation.
func NewPublicELBsDataSource() datasource.DataSource {
	return &publicELBsDataSource{}
}

// publicELBsDataSource is the data source implementation.
type publicELBsDataSource struct {
	client *hashicups.Client
}

// publicELBsDataSourceModel maps the data source schema data.
type publicELBsDataSourceModel struct {
	Arns types.List            `tfsdk:"arns"`
	Regions types.List         `tfsdk:"regions"`
}

func (p publicELBsDataSourceModel) GetRegions() types.List {
   return p.Regions
}


// Metadata returns the data source type name.
func (d *publicELBsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "public_elbs"
}

// Schema defines the schema for the data source.
func (d *publicELBsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
func (d *publicELBsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func addToV1Arns(
	ctx context.Context,
	elbV1RegionClient *elasticloadbalancing.Client,
	accountID string,
	region string,
	elbV1LoadBalancerARNs *[]string,
	nextMarker *string,
) (*string, error) {
	input := &elasticloadbalancing.DescribeLoadBalancersInput{
		Marker:   nextMarker,
	}

	// Call the DescribeLoadBalancers operation for ELBv1 in the current region
	resp, err := elbV1RegionClient.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, err
	}

	var lbARN string
	// Add 'internet-facing' LB ARNs
	for _, lb := range resp.LoadBalancerDescriptions {
		if *lb.Scheme == "internet-facing" {
			lbARN = fmt.Sprintf(
				"arn:aws:elasticloadbalancing:%s:%s:loadbalancer/%s",
				region,
				accountID,
				*lb.LoadBalancerName,
			)
			*elbV1LoadBalancerARNs = append(*elbV1LoadBalancerARNs, lbARN)
		}
	}

	return resp.NextMarker, nil
}

func addToV2Arns(
	ctx context.Context,
	elbV2RegionClient *elasticloadbalancingv2.Client,
	elbV2LoadBalancerARNs *[]string,
	nextMarker *string,
) (*string, error) {
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{
		Marker:   nextMarker,
	}

	// Call the DescribeLoadBalancers operation for ELBv2 in the current region
	resp, err := elbV2RegionClient.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, err
	}

	// Add 'internet-facing' LB ARNs
	for _, lb := range resp.LoadBalancers {
		if lb.Scheme ==  "internet-facing" {
			*elbV2LoadBalancerARNs = append(*elbV2LoadBalancerARNs, *lb.LoadBalancerArn)
		}
	}

	return resp.NextMarker, nil
}

// Read refreshes the Terraform state with the latest data.
func (d *publicELBsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state publicELBsDataSourceModel

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

	regions, err := getRegions(ctx, cfg, state, resp)
	if err != nil {
		return
	}

	accountID, err := getAccountID(ctx, cfg, resp)
	if err != nil {
		return
	}

	// Initialize variables for pagination
	var nextMarker *string
	var loadBalancerARNs []string


	// Iterate through regions
	for _, region := range *regions {
		// Create a client for the current region
		regionCfg := cfg.Copy()
		regionCfg.Region = region
		elbV1RegionClient := elasticloadbalancing.NewFromConfig(regionCfg)
		elbV2RegionClient := elasticloadbalancingv2.NewFromConfig(regionCfg)

		for {
			nextMarker, err = addToV1Arns(
				ctx,
				elbV1RegionClient,
				accountID,
				region,
				&loadBalancerARNs,
				nextMarker,
			)
			if err != nil {
				resp.Diagnostics.AddError("Error describing ELBv1 load balancers in region"+ region + ":" + err.Error(), "")
				return
			}
			// Check if there are more pages of ELBv1 load balancers to retrieve
			if nextMarker == nil {
				break
			}
		}

		for {
			nextMarker, err = addToV2Arns(ctx, elbV2RegionClient, &loadBalancerARNs, nextMarker)
			if err != nil {
				resp.Diagnostics.AddError("Error describing ELBv2 load balancers in region" + region + ":" + err.Error(), "")
				return
			}
			// Check if there are more pages of ELBv2 load balancers to retrieve
			if nextMarker == nil {
				break
			}
		}
	}
	// Sort lists
	sort.Strings(loadBalancerARNs)
	sort.Strings(*regions)

	// Ignore diags since we would catch error elsewhere
	state.Arns, _ = types.ListValueFrom(ctx, types.StringType, loadBalancerARNs)
	state.Regions, _ = types.ListValueFrom(ctx, types.StringType, regions)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
