package hashicups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type dataSourceCoffeesType struct{}

func (r dataSourceCoffeesType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{}, nil
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
}
