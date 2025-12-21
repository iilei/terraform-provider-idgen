#!/usr/bin/env bash
set -eo pipefail

# Usage:
# interactive:
# > ./preflight.sh
#   *  'seed_prefix'        # (may be blank; "")
#   *  'iteration_count'    # how many examples
#   *  'iteration_offset'   # offset iterator

# Non-interactive: pass positional arguments
# seeds 0..4999:
# > ./preflight.sh '' 5000 0

# seeds "x*-0".."x*-4999":
# > ./preflight.sh 'x*-' 5000 0

TEST_DIR="./preflight"
ID_GEN_VERSION="0.0.2"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

# Accept arguments or prompt for input
TF_VAR_string_prefix="${1:-}"
iteration_count="${2:-}"
iteration_offset="${3:-}"

# Accept arguments or prompt for input
if [ $# -ge 1 ]; then
  TF_VAR_string_prefix="$1"
else
  read -p "What to prepend to seed: " TF_VAR_string_prefix
fi

if [ $# -ge 2 ]; then
  iteration_count="$2"
else
  read -p "How many iterations: " iteration_count
fi

if [ $# -ge 3 ]; then
  iteration_offset="$3"
else
  read -p "Iteration offset: " iteration_offset
fi


# Generate main.tf with multiple data sources
cat > "$TEST_DIR/main.tf" <<EOF
terraform {
  required_providers {
    idgen = {
      source  = "iilei/idgen"
      version = "${ID_GEN_VERSION}"
    }
  }
}


variable "string_prefix" {
  type        = string
  description = "what to prepend to seed"
  default     = ""
}
EOF

# Add data sources for numbers 0-100000
for i in $(seq $iteration_offset $(($iteration_offset + $iteration_count - 1)) ); do
  cat >> "$TEST_DIR/main.tf" <<EOF

data "idgen_proquint" "test_${i}" {
  length = 11
  seed   = "${TF_VAR_string_prefix}${i}"
}
EOF
done

# add all generated ids to outputs
for i in $(seq $iteration_offset $(($iteration_offset + $iteration_count - 1)) ); do
  cat >> "$TEST_DIR/main.tf" <<EOF

output "idgen_proquint__test_${i}__id" {
   value = data.idgen_proquint.test_${i}.id
}
EOF
done

# Add outputs
cat >> "$TEST_DIR/main.tf" <<EOF
output "count" {
  value = "Generated $iteration_count proquints"
}
EOF


# Add outputs
cat >> "$TEST_DIR/main.tf" <<EOF
output "first_seed" {
  value = "${string_prefix}${iteration_offset}"
}
output "last_seed" {
  value = "${string_prefix}$(($iteration_offset + $iteration_count - 1))"
}
EOF

cd "$TEST_DIR" >/dev/null
terraform init >/dev/null

terraform apply -auto-approve >/dev/null

# Results as tsv
printf "Seed \tProquint"
echo ""
terraform output -json | jq -r '
  to_entries |
  map(select(.key | startswith("idgen_proquint__test_"))) |
  map({
    seed: (.key | capture("test_(?<num>[0-9]+)").num),
    id: .value.value
  }) |
  sort_by(.seed | tonumber) |
  .[] |
  "'${TF_VAR_string_prefix}'\(.seed) \t\(.id)"
'

# Show duplicates summary
echo ""
echo "# === Duplicate Analysis ==="
terraform output -json | jq -r '
  (
    to_entries |
    map(select(.key | startswith("idgen_proquint__test_"))) |
    map({
      seed: (.key | capture("test_(?<num>[0-9]+)").num),
      id: .value.value
    }) |
    group_by(.id) |
    map(select(length > 1))
  ) as $duplicates |
  if ($duplicates | length) == 0 then
    "# No duplicates found - all IDs are unique for '${iteration_count}' seeds."
  else
    $duplicates |
    map({
      id: .[0].id,
      count: length,
      seeds: map(.seed)
    }) |
    .[] |
    "# Seeds: \(.seeds | map("'"${TF_VAR_string_prefix}"'\(.)") | join(", ")) -- \(.count) times yields same ID: \(.id)"
  end
'
cd - >/dev/null

rm -rf "$TEST_DIR"
