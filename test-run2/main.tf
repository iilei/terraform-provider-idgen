terraform {
  required_providers {
    idgen = {
      source  = "localhost/local/idgen"
    }
  }
}

provider "idgen" {}

data "idgen_proquint" "min_value" {
  length     = 11
  seed       = 0
}

data "idgen_proquint" "max_value" {
  length     = 11
  seed       = 4294967295
}

data "idgen_proquint_canonical" "maximum_value" {
  seed       = "4294967296"
}




# Test Templated ID Generation, see:
# internal/provider/docs_embeds/templated_data_source.md

# Test Case Conversion Functions

# Input: "sunny" | Output: "SUNNY"
data "idgen_templated" "test_upper" {
  template = "{{ .random_word | upper }}"
  random_word = { seed = "17" }
}

# Input: "SUNNY" | Output: "sunny"
data "idgen_templated" "test_lower" {
  template = "{{ .random_word | upper | lower }}"
  random_word = { seed = "17" }
}

# Test String Manipulation Functions

# Input: "sunny" | Output: "sun"
data "idgen_templated" "test_replace" {
  template = "{{ .random_word | replace \"ny\" \"\" }}"
  random_word = { seed = "17" }
}

# Input: "sunny" | Output: "unn"
data "idgen_templated" "test_substr" {
  template = "{{ .random_word | substr 1 3 }}"
  random_word = { seed = "17" }
}

# Input: "  sunny  " | Output: "sunny"
data "idgen_templated" "test_trim" {
  template = "{{ .random_word | trim }}"
  random_word = { seed = "17" }
}

# Input: "sunny" | Output: "nny"
data "idgen_templated" "test_trimPrefix" {
  template = "{{ .random_word | trimPrefix \"su\" }}"
  random_word = { seed = "17" }
}

# Input: "sunny" | Output: "sun"
data "idgen_templated" "test_trimSuffix" {
  template = "{{ .random_word | trimSuffix \"ny\" }}"
  random_word = { seed = "17" }
}

# Test Repetition & Reversal Functions

# Input: "sunny" | Output: "sunnysunnysunny"
data "idgen_templated" "test_repeat" {
  template = "{{ .random_word | repeat 3 }}"
  random_word = { seed = "17" }
}

# Input: "sunny" | Output: "ynnus"
data "idgen_templated" "test_reverse" {
  template = "{{ .random_word | reverse }}"
  random_word = { seed = "17" }
}

# Test More Complex Examples

# yields: 0q-LUSAB_BABAD
data "idgen_templated" "example1" {
  template = "0q-{{ .proquint_canonical | upper | replace \"-\" \"_\" }}"
  proquint_canonical = { seed = "127.0.0.1" }
}

# yields: snowy-vibub-vamiz
data "idgen_templated" "example2" {
  template = "{{ .random_word }}-{{ .proquint }}"
  random_word = { seed = "a5e57e8a-9a7c-4efd-9fdd-0fcdc7630e3a" }
  proquint = { seed = "a5e57e8a-9a7c-4efd-9fdd-0fcdc7630e3a" }
}

# yields: bCi-pPV
data "idgen_templated" "example3" {
  template = "{{ .nanoid }}"
  nanoid = { length = 21, seed = "72da0233-3b03-4410-854f-3b96e868e15a", alphabet = "readable", length = 7, group_size = 3 }
}

# Output values to verify correctness

output "test_proquint_maximum_value" { value = data.idgen_proquint.maximum_value.id }
output "test_proquint_max" { value = data.idgen_proquint.max_value.id }
output "test_proquint_min" { value = data.idgen_proquint.min_value.id }
output "test_upper" { value = data.idgen_templated.test_upper.id }
output "test_lower" { value = data.idgen_templated.test_lower.id }
output "test_replace" { value = data.idgen_templated.test_replace.id }
output "test_substr" { value = data.idgen_templated.test_substr.id }
output "test_trim" { value = data.idgen_templated.test_trim.id }
output "test_trimPrefix" { value = data.idgen_templated.test_trimPrefix.id }
output "test_trimSuffix" { value = data.idgen_templated.test_trimSuffix.id }
output "test_repeat" { value = data.idgen_templated.test_repeat.id }
output "test_reverse" { value = data.idgen_templated.test_reverse.id }
output "example1" { value = data.idgen_templated.example1.id }
output "example2" { value = data.idgen_templated.example2.id }
output "example3" { value = data.idgen_templated.example3.id }
