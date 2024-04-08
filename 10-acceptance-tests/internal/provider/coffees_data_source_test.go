package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCoffeesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "hashicups_coffees" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.#", "9"),
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.description", ""),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.id", "1"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.image", "/hashicorp.png"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.ingredients.#", "1"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.ingredients.0.id", "6"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.name", "HCP Aeropress"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.price", "200"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.teaser", "Automation in a cup"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "id", "placeholder"),
				),
			},
		},
	})
}
