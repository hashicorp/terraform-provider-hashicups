---
page_title: "order Data Source - terraform-provider-hashicups"
subcategory: ""
description: |-
  The order data source allows you to retrieve information about a particular HashiCups order.
---

# Data Source `hashicups_order`

-> Visit the [Perform CRUD operations with Providers](https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers&utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS) Learn tutorial for an interactive getting started experience.

The order data source allows you to retrieve information about a particular HashiCups order.

## Example Usage

```terraform
data "hashicups_order" "edu" {
  id = 1
}
```

## Argument Reference

- `id` - (Required) HashiCups order ID.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `items` - Items in a HashiCups order. See [Order item](#order-item) below for details.

### Order item

Each order item contains a `coffee` object and a `quantity`.

- `coffee` - Represents a HashiCups coffee object. See [Coffee](#coffee) below for details.
- `quantity` - The number of coffee in an order item.

### Coffee

- `id` -  The coffee ID.
- `image` - The coffee's image URL path.
- `name` - The coffee name.
- `price` - The coffee price.
- `teaser` - The coffee teaser.
