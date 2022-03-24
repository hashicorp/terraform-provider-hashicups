#!/bin/bash

HTTP_RESPONSE=$(curl \
   --header "Authorization: Bearer $TFC_TOKEN" \
   --header "Content-Type: application/vnd.api+json" \
   --request POST \
   --data @public_key.json \
   https://app.terraform.io/api/registry/private/v2/gpg-keys | jq -r '.data | .attributes | ."key-id"')

echo $HTTP_RESPONSE