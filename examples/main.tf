terraform {
  required_providers {
    hashicups = {
      version = "0.3"
      source = "hashicorp.com/edu/hashicups"
    }
  }
}

provider "hashicups" {
  username = "dos"
  password = "test123"
}

module "psl" {
  source = "./coffee"

  coffee_name = "Packer Spiced Latte"
}

output "psl" {
  value = module.psl.coffee
}

data "hashicups_ingredients" "psl" {
  coffee_id = values(module.psl.coffee)[0].id
}

# output "psl_i" {
#   value = data.hashicups_ingredients.psl
# }

resource "hashicups_order" "new" {
  items {
    coffee {
      id = 3
    }
    quantity = 2
  }
  items {
    coffee {
      id = 2
    }
    quantity = 2
  }
}

output "new_order" {
  value = hashicups_order.new
}


data "hashicups_order" "first" {
  id = 1
}

output "first_order" {
  value = data.hashicups_order.first
}
