#!/usr/bin/env bash
set -eo pipefail

# Build the provider
echo "Building provider..."
go clean -cache
go build -o terraform-provider-idgen

# Create a test directory
TEST_DIR_1="./test-run1"
TEST_DIR_2="./test-run2"
TEST_DIR_3="./test-run3"
rm -rf "$TEST_DIR_1"
rm -rf "$TEST_DIR_2"
rm -rf "$TEST_DIR_3"
mkdir -p "$TEST_DIR_1"
mkdir -p "$TEST_DIR_2"
mkdir -p "$TEST_DIR_3"

# Define test seeds
SEED_PREFIX="${TF_VAR_seed_prefix:-asdf}"
SEEDS=(
  "$SEED_PREFIX"
  "${SEED_PREFIX}-1"
  "${SEED_PREFIX}-2"
  "${SEED_PREFIX}-3"
  "${SEED_PREFIX}-4"
  "${SEED_PREFIX}-5"
  "${SEED_PREFIX}-6"
  "${SEED_PREFIX}-7"
  "${SEED_PREFIX}-8"
  "${SEED_PREFIX}-9"
  "${SEED_PREFIX}-10"
  "${SEED_PREFIX}-11"
  "${SEED_PREFIX}-12"
)

# Copy provider binary to test directory with proper naming
PLUGIN_DIR="$TEST_DIR_1/.terraform/plugins/localhost/local/idgen/0.0.1/$(go env GOOS)_$(go env GOARCH)"
mkdir -p "$PLUGIN_DIR"
cp terraform-provider-idgen "$PLUGIN_DIR/"

# Create terraform CLI config to override provider
cat > "$TEST_DIR_1/.terraformrc" <<EOF
provider_installation {
  dev_overrides {
    "localhost/local/idgen" = "$PWD/$TEST_DIR_1/.terraform/plugins/localhost/local/idgen/0.0.1/$(go env GOOS)_$(go env GOARCH)"
  }
  direct {}
}
EOF

# Create test configuration
cat > "$TEST_DIR_1/main.tf" <<'EOF'
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

data "idgen_nanoid" "test" {
  length     = 7
  group_size = 3
  alphabet   = "alphanumeric"
  seed       = var.seed
}

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

EOF

cat > "$TEST_DIR_2/.terraformrc" <<EOF
provider_installation {
  dev_overrides {
    "localhost/local/idgen" = "$PWD/$TEST_DIR_1/.terraform/plugins/localhost/local/idgen/0.0.1/$(go env GOOS)_$(go env GOARCH)"
  }
  direct {}
}
EOF

cat > "$TEST_DIR_2/main.tf" <<'EOF'
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


data "idgen_templated" "test_upper" {
  template = "{{ .random_word | upper }}"
  random_word = { seed = "17" }
}

data "idgen_templated" "test_lower" {
  template = "{{ .random_word | upper | lower }}"
  random_word = { seed = "17" }
}


data "idgen_templated" "test_replace" {
  template = "{{ .random_word | replace \"ny\" \"\" }}"
  random_word = { seed = "17" }
}

data "idgen_templated" "test_substr" {
  template = "{{ .random_word | substr 1 3 }}"
  random_word = { seed = "17" }
}

data "idgen_templated" "test_trim" {
  template = "{{ .random_word | trim }}"
  random_word = { seed = "17" }
}

data "idgen_templated" "test_trimPrefix" {
  template = "{{ .random_word | trimPrefix \"su\" }}"
  random_word = { seed = "17" }
}

data "idgen_templated" "test_trimSuffix" {
  template = "{{ .random_word | trimSuffix \"ny\" }}"
  random_word = { seed = "17" }
}

# yields "vivid-vivid-vivid"
data "idgen_templated" "test_repeat" {
  template = "{{ .random_word | append \"-\" | repeat 3 | trimSuffix \"-\"}}"
  random_word = { seed = "17" }
}

data "idgen_templated" "test_reverse" {
  template = "{{ .random_word | reverse }}"
  random_word = { seed = "17" }
}


data "idgen_templated" "example1" {
  template = "0q-{{ .proquint_canonical | upper | replace \"-\" \"_\" }}"
  proquint_canonical = { seed = "127.0.0.1" }
}

data "idgen_templated" "example2" {
  template = "{{ .random_word }}-{{ .proquint }}"
  random_word = { seed = "a5e57e8a-9a7c-4efd-9fdd-0fcdc7630e3a" }
  proquint = { seed = "a5e57e8a-9a7c-4efd-9fdd-0fcdc7630e3a" }
}

data "idgen_templated" "example3" {
  template = "{{ .nanoid }}"
  nanoid = { length = 21, seed = "72da0233-3b03-4410-854f-3b96e868e15a", alphabet = "readable", length = 7, group_size = 3 }
}


locals {
   stage = "dev"
   seed = "app-specific-seed"
   size = 4
   size_fmt  = format("%03d", local.size)  # "004"
}

# result: "0q-zozif-zapuf-rXK-s004dev"
data "idgen_templated" "infra_naming_docs_example" {
  template = "0q-{{ .proquint }}-{{ .nanoid }}-s${local.size_fmt}-${local.stage}"
  nanoid = { length = 3, seed = "#${local.size}_${local.seed}", alphabet = "readable" }
  proquint = { length = 11, seed = "#${local.size}_${local.seed}" }
}

# Output values to verify correctness

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
output "infra_naming_docs_example" { value = data.idgen_templated.infra_naming_docs_example.id }
EOF


# Create terraform CLI config to override provider
cat > "$TEST_DIR_2/.terraformrc" <<EOF
provider_installation {
  dev_overrides {
    "localhost/local/idgen" = "$PWD/$TEST_DIR_.terraform/plugins/localhost/local/idgen/0.0.1/$(go env GOOS)_$(go env GOARCH)"
  }
  direct {}
}
EOF

# Run terraform
cd "$TEST_DIR_1"
export TF_CLI_CONFIG_FILE="$PWD/.terraformrc"

echo ""
echo "Testing ID generation with different seeds"
echo "=========================================="
printf "%-30s | %-15s | %s\n" "Seed" "NanoID"  "Proquint"
echo "------------------------------------------------------------------------"

for SEED in "${SEEDS[@]}"; do
  # Run terraform and capture output
  FULL_OUTPUT=$(terraform apply -auto-approve -var="seed=$SEED" 2>&1)

  # Extract the last occurrence of nanoid and proquint values
  NANOID=$(echo "$FULL_OUTPUT" | grep 'nanoid = ' | tail -1 | sed 's/.*= "\(.*\)"/\1/')
  PROQUINT=$(echo "$FULL_OUTPUT" | grep 'proquint = ' | tail -1 | sed 's/.*= "\(.*\)"/\1/')

  printf "%-30s | %-15s | %s\n" "$SEED" "$NANOID" "$PROQUINT"
done


cd -

cd "$TEST_DIR_2"

terraform apply -auto-approve

cd -

# Test templated-parametrization example
echo ""
echo "Testing templated-parametrization example"
echo "========================================="

# Set up TEST_DIR_3 for templated-parametrization
PLUGIN_DIR_3="$TEST_DIR_3/.terraform/plugins/localhost/local/idgen/0.0.1/$(go env GOOS)_$(go env GOARCH)"
mkdir -p "$PLUGIN_DIR_3"
cp terraform-provider-idgen "$PLUGIN_DIR_3/"

cat > "$TEST_DIR_3/.terraformrc" <<EOF
provider_installation {
  dev_overrides {
    "localhost/local/idgen" = "$PWD/$TEST_DIR_3/.terraform/plugins/localhost/local/idgen/0.0.1/$(go env GOOS)_$(go env GOARCH)"
  }
  direct {}
}
EOF

# Copy the templated-parametrization example and modify provider source for local testing
cp examples/templated-parametrization/main.tf "$TEST_DIR_3/"
cp examples/templated-parametrization/variables.tf "$TEST_DIR_3/"

# Override provider source to use local version for testing
sed -i 's|source = "registry.terraform.io/iilei/idgen"|source = "localhost/local/idgen"|' "$TEST_DIR_3/main.tf"

cd "$TEST_DIR_3"
export TF_CLI_CONFIG_FILE="$PWD/.terraformrc"

echo "Running with default values:"
terraform apply -auto-approve

echo ""
echo "Running with custom values:"
terraform apply -auto-approve \
  -var="app_seed=seed-by-team-xyz" \
  -var="environment=dev" \
  -var="cluster_size=4" \
  -var="app_version=7" \
  -var="app_name=xyz" \
  -var="region=eu_central_1"

cd -

rm -rf "$TEST_DIR_3"
rm -rf "$TEST_DIR_2"
rm -rf "$TEST_DIR_1"

echo ""
echo "=========================================="
echo "Try different seed prefixes:"
echo "  TF_VAR_seed_prefix=myapp ./test-provider.sh"
echo "  TF_VAR_seed_prefix="127.0.0.1" ./test-provider.sh"
echo '  TF_VAR_seed_prefix=$(date +"x%s%N") ./test-provider.sh'
echo "=========================================="
