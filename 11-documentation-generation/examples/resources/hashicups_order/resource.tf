# Copyright IBM Corp. 2021, 2026

# Manage example order.
resource "hashicups_order" "example" {
  items = [
    {
      coffee = {
        id = 3
      }
      quantity = 2
    },
  ]
}
