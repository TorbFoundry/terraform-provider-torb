terraform {
  required_providers {
    torb = {
      version = "0.2"
      source  = "torb/tf/provider"
    }
  }
}

provider "torb" {}

module "helm_release" {
  source = "./data-sources/helm_release"
  release_name = "hello-world"
  namespace = "terraform-provider-torb-testing"
}

output "all_values" {
  value = jsondecode(module.helm_release.all_values)
}