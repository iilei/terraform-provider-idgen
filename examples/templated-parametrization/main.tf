terraform {
  required_providers {
    idgen = {
      source = "registry.terraform.io/iilei/idgen"
    }
  }
}

provider "idgen" {}

# Generate a single templated ID combining Proquint, NanoID and local variables.
# Seeds result in fully deterministic IDs.
locals {
  # The seed is going to produce deterministic IDs for both nanoid and proquint
  seed = var.app_seed

  # Other variables for parametrization
  size     = var.cluster_size
  size_fmt = format("%03d", local.size)  # zero-padded 3 digits
  stage    = var.environment
}

data "idgen_templated" "example" {
  template = "0q-{{ .proquint }}-{{ .nanoid }}-s${local.size_fmt}${local.stage}"
  nanoid   = { length = 3, seed = "#${local.size}_${local.seed}", alphabet = "readable" }
  proquint = { length = 11, seed = "#${local.size}_${local.seed}" }
}

# Multiple variations showing different parametrization patterns
data "idgen_templated" "infrastructure_name" {
  template = "${var.environment}-{{ .proquint | upper }}-cluster-{{ .nanoid }}"
  nanoid   = { length = 6, seed = "${var.app_seed}-infra", alphabet = "alphanumeric", group_size = 3 }
  proquint = { length = 11, seed = "${var.app_seed}-infra", group_size = 5 }
}

data "idgen_templated" "versioned_resource" {
  template = "{{ .random_word }}-v${format("%02d", var.app_version)}-{{ .nanoid }}"
  random_word = { seed = "${var.app_seed}-resource" }
  nanoid      = { length = 8, seed = "${var.app_seed}-v${var.app_version}", alphabet = "readable" }
}

# Database naming with compliance formatting
data "idgen_templated" "database_name" {
  template = "{{ .random_word | lower }}_${lower(var.environment)}_{{ .nanoid | lower }}"
  random_word = { seed = "${var.app_seed}-db" }
  nanoid      = { length = 8, seed = "${var.app_seed}-db", alphabet = "alphanumeric" }
}

# S3 bucket with region and stage
data "idgen_templated" "s3_bucket" {
  template = "${lower(var.environment)}-${replace(var.region, "_", "-")}-0q-{{ .proquint }}"
  proquint = { seed = "${var.app_seed}-storage" }
}

output "templated_id_basic" {
  value       = data.idgen_templated.example.id
  description = "Basic templated ID with parametrization"
}

output "infrastructure_name" {
  value       = data.idgen_templated.infrastructure_name.id
  description = "Infrastructure resource name with uppercase proquint"
}

output "versioned_resource_name" {
  value       = data.idgen_templated.versioned_resource.id
  description = "Versioned resource name with random word and version number"
}

output "database_identifier" {
  value       = data.idgen_templated.database_name.id
  description = "Database name following lowercase underscore convention"
}

output "s3_bucket_name" {
  value       = data.idgen_templated.s3_bucket.id
  description = "S3 bucket name compliant with AWS naming conventions"
}

# Show all parameters used
output "parameters_used" {
  value = {
    app_seed     = var.app_seed
    environment  = var.environment
    cluster_size = var.cluster_size
    app_version  = var.app_version
    app_name     = var.app_name
    region       = var.region
    size_formatted = local.size_fmt
  }
  description = "All input parameters used for templating"
}

# Expected outputs with test values:
# templated_id_basic = "0q-nivis-zozak-QDJ-s004dev"
# infrastructure_name = "dev-LADOZ-ZABAJ-cluster-9W7-kZ"
# versioned_resource_name = "minty-v07-CngFs5K6"
# database_identifier = "rural_dev_86qo8vnu"
# s3_bucket_name = "dev-eu-central-1-0q-nomil-tiput"
