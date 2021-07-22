---
page_title: "coffee Data Source - terraform-provider-hashicups"
subcategory: ""
description: |-
  The coffee data source allows you to retrieve information all available HashiCups coffees.
---

# Data Source `coffee`

The coffee data source allows you to retrieve information all available HashiCups coffees.

## Example Usage

```terraform
data "hashicups_coffees" "all" {}

```

## Attributes Reference

The following attributes are exported.

- `coffees` - A list of HashiCups coffee objects. See [Coffee](#coffee) below for details.

### Coffee

- `id` -  The coffee ID.
- `image` - The coffee's image URL path.
- `name` - The coffee name.
- `price` - The coffee price.
- `teaser` - The coffee teaser.
- `description` - The coffee description.
- `ingredients` - A list of coffee ingredients. See [Ingredients](#ingredients) below for details.

### Ingredients

- `ingredient_id` - The ingredient ID.