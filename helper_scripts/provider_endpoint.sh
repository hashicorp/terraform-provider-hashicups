#!/bin/bash

HTTP_RESPONSE=$(curl \
  --header "Authorization: Bearer $TFC_TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  --request POST \
  --data @provider.json \
  https://app.terraform.io/api/v2/organizations/<YOUR_ORG_NAME>/registry-providers)

echo $HTTP_RESPONSE