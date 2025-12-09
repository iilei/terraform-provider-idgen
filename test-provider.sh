#!/usr/bin/env bash
set -eo pipefail

# Build the provider
echo "Building provider..."
go build -o terraform-provider-idgen

# Create a test directory
TEST_DIR="./test-run"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

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
PLUGIN_DIR="$TEST_DIR/.terraform/plugins/localhost/local/idgen/0.0.1/$(go env GOOS)_$(go env GOARCH)"
mkdir -p "$PLUGIN_DIR"
cp terraform-provider-idgen "$PLUGIN_DIR/"

# Create terraform CLI config to override provider
cat > "$TEST_DIR/.terraformrc" <<EOF
provider_installation {
  dev_overrides {
    "localhost/local/idgen" = "$PWD/$TEST_DIR/.terraform/plugins/localhost/local/idgen/0.0.1/$(go env GOOS)_$(go env GOARCH)"
  }
  direct {}
}
EOF

# Create test configuration
cat > "$TEST_DIR/main.tf" <<'EOF'
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
EOF

# Run terraform
cd "$TEST_DIR"
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
rm -rf "$TEST_DIR"

echo ""
echo "=========================================="
echo "Try different seed prefixes:"
echo "  TF_VAR_seed_prefix=myapp ./test-provider.sh"
echo "  TF_VAR_seed_prefix="127.0.0.1" ./test-provider.sh"
echo '  TF_VAR_seed_prefix=$(date +"x%s%N") ./test-provider.sh'
echo "=========================================="
