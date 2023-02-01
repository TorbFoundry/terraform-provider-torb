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

module "torb_helm_release_test" {
  source       = "./data-sources/helm_release"
  release_name = "hello-world"
  namespace    = "test-torb"
}

output "all_values" {
  value = jsondecode(module.torb_helm_release_test.all_values)
}