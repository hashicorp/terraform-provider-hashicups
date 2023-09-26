terraform {
  required_providers {
    hashicups = {
      source = "hashicorp.com/edu/hashicups-pf"
    }
  }
}

provider "hashicups" {
  host     = "http://localhost:19090"
  username = "education"
  password = "test123"
}

data "hashicups_coffees" "edu" {
  regions = ["us-west-1",  "us-east-1", "us-east-2", "us-west-2"]
}

output "edu_coffees" {
  value = data.hashicups_coffees.edu
}