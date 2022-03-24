​​#!/bin/bash

HTTP_RESPONSE=$(curl \
  --header "Authorization: Bearer $TFC_TOKEN" \
  --header "Content-Type: text/html" \
  --request POST \
  --data  https://github.com/tr0njavolta/terraform-provider-hashicups-pf/releases/download/v0.1.0/terraform-provider-hashicups-pf_0.1.0_linux_amd64.zip \
    <PLATFORM_URL>)
)

echo $HTTP_RESPONSE
