terraform {
  required_providers {
    idgen = {
      source  = "iilei/idgen"
      version = "1.1.0"
    }
  }
}

# Encode a SHA-256 digest as a 95-character proquint (32 bytes, 16 words).
# sha256() returns lowercase hex without a prefix, which is the default seed_format.
data "idgen_proquint_big" "sha256_abc" {
  length = 95
  seed   = sha256("abc")
}

# The same digest written as a 0x-prefixed hex literal gives the same ID.
data "idgen_proquint_big" "sha256_abc_prefixed" {
  length = 95
  seed   = "0xba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
}

# Decimal input must be requested explicitly, because digit-only strings are valid hex too.
data "idgen_proquint_big" "sha256_abc_decimal" {
  length      = 95
  seed_format = "decimal"
  seed        = "84342368487090800366523834928142263660104883695016514377462985829716817089965"
}

output "sha256_abc" {
  value = data.idgen_proquint_big.sha256_abc.id
}

output "sha256_abc_prefixed" {
  value = data.idgen_proquint_big.sha256_abc_prefixed.id
}

output "sha256_abc_decimal" {
  value = data.idgen_proquint_big.sha256_abc_decimal.id
}
