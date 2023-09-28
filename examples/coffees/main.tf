terraform {
  required_providers {
    public = {
      source = "hashicorp.com/edu/hashicups-pf"
    }
  }
}

provider "public" {
  host     = "http://localhost:19090"
  username = "education"
  password = "test123"
}

data "public_subnets" "subnets" {
  regions = ["us-west-1",  "us-east-1", "us-east-2", "us-west-2"]
}

output "public_subnets" {
  value = data.public_subnets.subnets
}

data "public_ec2s" "ec2s" {
  regions = ["us-west-1",  "us-east-1", "us-east-2", "us-west-2"]
}

output "public_ec2s" {
  value = data.public_ec2s.ec2s
}