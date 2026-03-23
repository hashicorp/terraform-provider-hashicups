terraform {
  required_providers {
    hashicups = {
      source  = "localhost:5601/hashicups/hashicups"
      version = "0.0.1"
    }
  }
}

provider "hashicups" {
  host     = "http://localhost:19090"
  username = "education"
  password = "test123"
}

data "hashicups_coffees" "edu" {}

output "edu_coffees" {
  value = data.hashicups_coffees.edu
}
