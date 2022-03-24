#!/bin/bash

HTTP_RESPONSE=$(curl \
  --header "Authorization: Bearer $TFC_TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  --request POST \
  --data @version.json \
  https://app.terraform.io/api/v2/organizations/<YOUR_ORG_NAME>/registry-providers/private/<YOUR_ORG_NAME>/hashicups/versions)

SHA_RESPONSE=$(curl \
  --header "Authorization: Bearer $TFC_TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  "https://app.terraform.io/api/v2/organizations/<YOUR_ORG_NAME>/registry-providers/private/<YOUR_ORG_NAME>/hashicups/versions" | jq -r '.data[] | .links | ."shasums-upload"')

SIG_RESPONSE=$(curl \
  --header "Authorization: Bearer $TFC_TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  "https://app.terraform.io/api/v2/organizations/<YOUR_ORG_NAME>/registry-providers/private/<YOUR_ORG_NAME>/hashicups/versions" | jq -r '.data[] | .links | ."shasums-sig-upload"')

echo "shasums-upload-url:" $SHA_RESPONSE
echo "_"
echo "shasums-sig-upload-url:" $SIG_RESPONSE
