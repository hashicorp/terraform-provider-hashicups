package hashicups

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.#", "6"),
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.description", ""),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.id", "1"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.image", "/packer.png"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.ingredients.#", "3"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.ingredients.0.id", "1"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.ingredients.1.id", "2"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.ingredients.2.id", "4"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.name", "Packer Spiced Latte"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.price", "350"),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "coffees.0.teaser", "Packed with goodness to spice up your images"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "id", "placeholder"),
				),
			},
		},
	})
}
