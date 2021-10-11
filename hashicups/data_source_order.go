package hashicups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type dataSourceOrderType struct{}

func (r dataSourceOrderType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{}, nil
}

func (r dataSourceOrderType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceOrder{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceOrder struct {
	p provider
}

func (r dataSourceOrder) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
}
