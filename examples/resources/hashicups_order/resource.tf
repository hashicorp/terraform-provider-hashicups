# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

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