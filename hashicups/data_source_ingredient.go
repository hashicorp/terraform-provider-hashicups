package hashicups

import (
	"context"
	"fmt"
	"strconv"

	hc "github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIngredients() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIngredientsRead,
		Schema: map[string]*schema.Schema{
			"coffee_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ingredients": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIngredientsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hc.Client)

	// // Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	coffeeID := d.Get("coffee_id").(int)
	cID := strconv.Itoa(coffeeID)

	ings, err := c.GetCoffeeIngredients(cID)
	if err != nil {
		return diag.FromErr(err)
	}

	ingredients := make([]map[string]interface{}, 0)

	for _, v := range ings {
		ingredient := make(map[string]interface{})

		ingredient["id"] = v.ID
		ingredient["name"] = fmt.Sprintf("ingredient - %+v", v.Name)
		ingredient["quantity"] = v.Quantity
		ingredient["unit"] = v.Unit

		ingredients = append(ingredients, ingredient)
	}

	if err := d.Set("ingredients", ingredients); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cID)

	return diags
}
