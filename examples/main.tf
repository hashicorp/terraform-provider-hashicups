terraform {
  required_providers {
    hashicups = {
      version = "~> 0.3.1"
      source  = "hashicorp.com/edu/hashicups"
    }
  }
}
provider "hashicups" {
  username = "rachel"
  password = "test123"
  host     = "http://localhost:19090"
}

// resource "hashicups_order" "edu" {
//   items = {
//     coffee = {
//       id = 3
//     }
//     quantity = 2
//   }
// }

resource "hashicups_order" "edu" {
  items = [{
    coffee = {
      id = 3
    }
    quantity = 2
    }, {
    coffee = {
      id = 1
    }
    quantity = 2
    }
  ]
}

output "edu_order" {
  value = hashicups_order.edu
}
