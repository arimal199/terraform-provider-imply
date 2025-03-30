terraform {
  required_providers {
    imply = {
      source = "registry.terraform.io/arimal199/imply" // registry.opentofu.org/arimal199/imply
      # version = "0.0.1" // TODO: add version
    }
  }
}

provider "imply" {
  # host = "" // or use the IMPLY_HOST environment variable
  # api_key = "" // or use the IMPLY_API_KEY environment variable
}

data "imply_users" "_" {}

# data "imply_user" "_" {} // TODO: add user data source

data "imply_groups" "_" {}

data "imply_permissions" "_" {}
