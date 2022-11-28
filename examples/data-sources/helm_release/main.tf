terraform {
  # required_providers {
  #   torb = {
  #     version = "0.2"
  #     source  = "torb/tf/provider"
  #   }
  # }
}

variable "release_name" {
  type    = string
  default = "hello-world"
}

variable "namespace" {
  type    = string
  default = "torb-test"
}

# data "torb_helm_release" "release" {
#   release_name = var.release_name
#   namespace    = var.namespace
# }

# output "all_values" {
#   value = data.torb_helm_release.release.values
# }
