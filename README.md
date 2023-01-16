# Terraform Provider Hashicups

Demo of the computed nested attribute issue: 

if a schema contains a `Computed` nested attribute and one of the attibutes of the nested attribute is also `Computed` (e.g. `coffee` in `order` resource - it has few `Computed` attribues, e.g. `name`),
then `terraform plan` is always dirty, even if the configuration hasn't changed.

The repo is a modified version of [Terraform Hashicups example](https://github.com/hashicorp/terraform-provider-hashicups-pf).
The only diff from the vanilla `hashicups` example is setting `Computed` and `Optional` fields of `coffee` attribute in `order` [resource](hashicups/order_resource.go#L91)

## Reproducing the issue

The following command reproduces the error - it should fail with something like `order_resource_test.go:10: Step 1/3 error: After applying this test step, the plan was not empty`

```shell
make error
```
