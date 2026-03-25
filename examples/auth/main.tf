terraform {
  required_providers {
    imply = {
      source = "registry.opentofu.org/arimal199/imply"
      # version = "0.0.1" // TODO: add version
    }
  }
}

provider "imply" {
  # host = "" // or use the IMPLY_HOST environment variable
  # api_key = "" // or use the IMPLY_API_KEY environment variable
}

/* data "imply_users" "_" {}

data "imply_user" "_" {
  id = "16505d53-14c5-433a-84ca-00bbb9a2ae21"
}

data "imply_groups" "_" {}

data "imply_group" "_" {
  id = "b3b28dce-ac2a-4e5f-a840-0641a647a737"
}

data "imply_permissions" "_" {}
 */

resource "imply_user" "user" {
  username = "foo2@bar.com"
}
