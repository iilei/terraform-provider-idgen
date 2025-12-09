terraform {
  required_providers {
    idgen = {
      source = "registry.terraform.io/iilei/idgen"
    }
  }
}

provider "idgen" {}

data "idgen_nanoid" "example" {}

data "idgen_proquint" "example" {}

data "idgen_templated" "example" {}

output "nanoid" {
  value = data.idgen_nanoid.example.id
}

output "proquint" {
  value = data.idgen_proquint.example.id
}

output "templated" {
  value = data.idgen_templated.example.id
}
