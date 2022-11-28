terraform {
  required_providers {
    torb = {
      version = "0.2"
      source  = "torb/tf/provider"
    }
  }
}

provider "torb" {}

module "torb_helm_release_test" {
  source = "./data-sources/helm_release"
  release_name = "hello-world"
  namespace = "test-torb"
}

output "all_values" {
  value = jsondecode(module.torb_helm_release_test.all_values)
}