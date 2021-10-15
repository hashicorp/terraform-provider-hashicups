terraform {
  required_providers {
    hashicups = {
      source  = "hashicorp.com/edu/hashicups-pf"
    }
  }
}
provider "hashicups" {
  username = "education"
  password = "test123"
  host     = "http://localhost:19090"
}

resource "hashicups_order" "edu" {
  items = [{
    coffee = {
      id = 1
    }
    quantity = 2
    }, {
    coffee = {
      id = 1
    }
    quantity = 4
    }
  ]
}

output "edu_order" {
  value = hashicups_order.edu
}