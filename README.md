# Terraform Provider Hashicups

Demo of the computed nested attribute issue: 

`terraform plan` can show unempty plan for unchanged configurations with computed nested attributes that also have computed attributes.

The repo is a modified version of [Terraform Hashicups example](https://github.com/hashicorp/terraform-provider-hashicups-pf).

## Reproducing the issue

```shell
make error
```
