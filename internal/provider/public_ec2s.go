package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &publicEC2sDataSource{}
	_ datasource.DataSourceWithConfigure = &publicEC2sDataSource{}
)

// NewPublicEC2sDataSource is a helper function to simplify the provider implementation.
func NewPublicEC2sDataSource() datasource.DataSource {
	return &publicEC2sDataSource{}
}

// publicEC2sDataSource is the data source implementation.
type publicEC2sDataSource struct {
	client *hashicups.Client
}

// publicEC2sDataSourceModel maps the data source schema data.
type publicEC2sDataSourceModel struct {
	Arns types.List            `tfsdk:"arns"`
	Regions types.List         `tfsdk:"regions"`
}

func (p publicEC2sDataSourceModel) GetRegions() types.List {
   return p.Regions
}


// Metadata returns the data source type name.
func (d *publicEC2sDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "public_ec2s"
}

// Schema defines the schema for the data source.
func (d *publicEC2sDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
func (d *publicEC2sDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func getAccountID(ctx context.Context, cfg aws.Config, resp *datasource.ReadResponse) (string, error) {
	// Create an STS (Security Token Service) client
	client := sts.NewFromConfig(cfg)

	// Call the GetCallerIdentity operation to retrieve the current account ID
	getCallerIdentityResp, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		resp.Diagnostics.AddError("Error calling sts:GetCallerIdentity, " + err.Error(), "")
		return "", err
	}

	// Return the account ID as a string
	return *getCallerIdentityResp.Account, nil
}


func addToEC2Arns(
	ctx context.Context,
	ec2RegionClient *ec2.Client,
	accountID string,
	region string,
	ec2ARNs *[]string,
	nextToken *string,
) (*string, error) {
	input := &ec2.DescribeInstancesInput{
		NextToken:   nextToken,
	}

	// Call the DescribeInstances operation for EC2 in the current region
    resp, err := ec2RegionClient.DescribeInstances(ctx, input)
    if err != nil {
        return nil, err
    }

	var ec2ARN string

    for _, reservation := range resp.Reservations {
        for _, instance := range reservation.Instances {
            // Note: If a private instance gets associated with an EIP
            // Then, it will show up as 'PublicIpAddress'
            if instance.PublicIpAddress == nil {
                continue
            }
            ec2ARN = fmt.Sprintf(
            "arn:aws:ec2:%s:%s:instance/%s",
            region,
            accountID,
             *instance.InstanceId,
            )
            *ec2ARNs = append(*ec2ARNs, ec2ARN)
        }
    }

	return resp.NextToken, nil
}


// Read refreshes the Terraform state with the latest data.
func (d *publicEC2sDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state publicEC2sDataSourceModel

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
	var nextToken *string
	var ec2ARNs []string

	// Iterate through regions
	for _, region := range *regions {
		// Create a client for the current region
		regionCfg := cfg.Copy()
		regionCfg.Region = region
		regionClient := ec2.NewFromConfig(regionCfg)

		// Iterate through pages
		for {
			nextToken, err = addToEC2Arns(
				ctx,
				regionClient,
				accountID,
				region,
				&ec2ARNs,
				nextToken,
			)

			if err != nil {
				resp.Diagnostics.AddError("Error describing EC2 instances in region" + region + ": " + err.Error(), "")
				return
			}

			// Check if there are more pages of subnets to retrieve
			if nextToken == nil {
				break
			}
		}
	}
	// Sort lists
	sort.Strings(ec2ARNs)
	sort.Strings(*regions)

	// Ignore diags since we would catch error elsewhere
	state.Arns, _ = types.ListValueFrom(ctx, types.StringType, ec2ARNs)
	state.Regions, _ = types.ListValueFrom(ctx, types.StringType, regions)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
