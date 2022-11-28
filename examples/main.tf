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
  release_name = "dollar-bedroom"
  namespace = "flask-app-w-react-frontend"
}

output "all_values" {
  value = jsondecode(module.torb_helm_release_test.all_values)
}