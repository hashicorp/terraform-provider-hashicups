# Compute total price with tax
output "total_price" {
  value = provider::hashicups::compute_tax(5.00, 0.085)
}
