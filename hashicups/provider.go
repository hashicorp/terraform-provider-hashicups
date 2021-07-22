package hashicups

import (
	"context"
	"os"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var stderr = os.Stderr

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *hashicups.Client
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"username": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	Username types.String `tfsdk:"username"`
	Host     types.String `tfsdk:"host"`
	Password types.String `tfsdk:"password"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	err := req.Config.Get(ctx, &config)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error parsing configuration",
			Detail:   "Error parsing the configuration, this is an error in the provider. Please report the following to the provider developer:\n\n" + err.Error(),
		})
		return
	}

	// User must provide a user to the provider
	var username string
	if config.Username.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityWarning,
			Summary:  "Unable to create client",
			Detail:   "Cannot use unknown value as username",
		})
		return
	}

	if config.Username.Null {
		username = os.Getenv("HASHICUPS_USERNAME")
	} else {
		username = config.Username.Value
	}

	if username == "" {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			// Error vs warning - empty value must stop execution
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to find username",
			Detail:   "Username cannot be an empty string",
		})
	}

	// User must provide a password to the provider
	var password string
	if config.Password.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityWarning,
			Summary:  "Unable to create client",
			Detail:   "Cannot use unknown value as password",
		})
		return
	}

	if config.Password.Null {
		password = os.Getenv("HASHICUPS_PASSWORD")
	} else {
		password = config.Password.Value
	}

	if password == "" {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			// Error vs warning - empty value must stop execution
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to find password",
			Detail:   "password cannot be an empty string",
		})
	}

	// User must specify a host
	var host string
	if config.Host.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityWarning,
			Summary:  "Unable to create client",
			Detail:   "Cannot use unknown value as host",
		})
		return
	}

	if config.Host.Null {
		host = os.Getenv("HASHICUPS_HOST")
	} else {
		host = config.Host.Value
	}

	if host == "" {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			// Error vs warning - empty value must stop execution
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to find host",
			Detail:   "Host cannot be an empty string",
		})
	}

	// Create a new HashiCups client and set it to the provider client
	c, err := hashicups.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to create client",
			Detail:   "Unable to create hashicups client:\n\n" + err.Error(),
		})
		return
	}

	p.client = c
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.ResourceType{
		"hashicups_order": resourceOrderType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.DataSourceType{
		"hashicups_coffees":     dataSourceCoffeesType{},
	}, nil
}
