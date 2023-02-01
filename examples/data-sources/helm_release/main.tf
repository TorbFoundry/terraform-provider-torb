terraform {
  required_providers {
    torb = {
      source = "TorbFoundry/torb"
      version = "0.1.2"
    }
  }
}

provider "torb" {
  # Configuration options
}

variable "release_name" {
  type    = string
  default = "hello-world"
}

variable "namespace" {
  type    = string
  default = "torb-test"
}

data "torb_helm_release" "release" {
  release_name = var.release_name
  namespace    = var.namespace
}

output "all_values" {
  value = data.torb_helm_release.release.values
}
