package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"math"
)

// Ensure the implementation satisfies the desired interfaces.
var _ function.Function = &ComputeTaxFunction{}

type ComputeTaxFunction struct{}

func NewComputeTaxFunction() function.Function {
	return &ComputeTaxFunction{}
}

func (f *ComputeTaxFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "compute_tax"
}

func (f *ComputeTaxFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Compute tax for coffee",
		Description: "Given a price and tax rate, return the total cost including tax.",

		Parameters: []function.Parameter{
			function.Float64Parameter{
				Name:        "price",
				Description: "Price of coffee item.",
			},
			function.Float64Parameter{
				Name:        "rate",
				Description: "Tax rate. 0.085 == 8.5%",
			},
		},
		Return: function.Float64Return{},
	}
}

func (f *ComputeTaxFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var price float64
	var rate float64
	var total float64

	// Read Terraform argument data into the variables
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &price, &rate))

	total = math.Round((price+price*rate)*100) / 100

	// Set the result
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, total))
}
