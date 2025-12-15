terraform {
  required_providers {
    idgen = {
      source  = "localhost/local/idgen"
    }
  }
}

provider "idgen" {}

variable "seed" {
  type = string
}

# Test NanoID
data "idgen_nanoid" "test" {
  length     = 9
  group_size = 3
  alphabet   = "alphanumeric"
  seed       = var.seed
}

# Test Proquint
data "idgen_proquint" "test" {
  length     = 11
  group_size = 5
  seed       = var.seed
}

output "nanoid" {
  value = data.idgen_nanoid.test.id
}

output "proquint" {
  value = data.idgen_proquint.test.id
}
