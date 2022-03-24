#!/bin/bash

HTTP_RESPONSE=$(curl \

  --header "Authorization: Bearer $TFC_TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  --request POST \
  --data @platforms.json \
  https://app.terraform.io/api/v2/organizations/<YOUR_ORG_NAME>/registry-providers/private/<YOUR_ORG_NAME>/hashicups/versions/0.1.0/platforms ​​| jq -r '.data[] | .links | ."provider-binary-upload"')

echo $HTTP_RESPONSE