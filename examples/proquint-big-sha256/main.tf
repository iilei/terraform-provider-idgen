terraform {
  required_providers {
    idgen = {
      source  = "iilei/idgen"
      version = "0.0.3"
    }
  }
}

# SHA-256 max value (32 bytes, all bits set)
# Both seeds below are equivalent hexadecimal representations.
data "idgen_proquint_big" "sha256_max_without_prefix" {
  length = 95
  seed   = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
}

data "idgen_proquint_big" "sha256_max_with_prefix" {
  length = 95
  seed   = "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
}

output "sha256_max_without_prefix" {
  value = data.idgen_proquint_big.sha256_max_without_prefix.id
}

output "sha256_max_with_prefix" {
  value = data.idgen_proquint_big.sha256_max_with_prefix.id
}
