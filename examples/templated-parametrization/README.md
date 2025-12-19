# Templated IDs with Parametrization

This example demonstrates how to create sophisticated templated identifiers that combine multiple ID types with Terraform variables and local values for complex parametrization patterns.

## Key Concepts

### Zero-Padded Formatting
Uses Terraform's `format()` function to create zero-padded numbers:
```hcl
size_fmt = format("%03d", local.size)  # 4 → "004"
```

### Deterministic Seeding
Creates predictable IDs by combining variables in seeds:
```hcl
seed = "#${local.size}_${local.seed}"
```

### Template Functions
Leverages string manipulation functions for compliance:
- `upper` / `lower` - Case conversion for naming conventions
- `replace` - Character substitution (underscores to hyphens)
- Template interpolation with `${var.name}` syntax

## Examples Included

### 1. Basic Parametrized Template
```
Template: "0q-{{ .proquint }}-{{ .nanoid }}-s${size_fmt}${stage}"
Output:   "0q-nivis-zozak-QDJ-s004dev"
```

### 2. Infrastructure Resource Naming
```
Template: "${environment}-{{ .proquint | upper }}-cluster-{{ .nanoid }}"
Output:   "dev-LADOZ-ZABAJ-cluster-9W7-kZ"
```

### 3. Versioned Resources
```
Template: "{{ .random_word }}-v${format("%02d", app_version)}-{{ .nanoid }}"
Output:   "minty-v07-CngFs5K6"
```

### 4. Database Naming (Underscore Convention)
```
Template: "{{ .random_word | lower }}_${lower(environment)}_{{ .nanoid | lower }}"
Output:   "rural_dev_86qo8vnu"
```

### 5. S3 Bucket Naming (AWS Compliant)
```
Template: "${lower(app_name)}-${lower(environment)}-${replace(region, "_", "-")}-{{ .proquint | lower }}"
Output:   "xyz-dev-eu-central-1-nomil-tiput"
```

## Expected Outputs

With the tested variable values:

| Output Name | Value | Use Case |
|-------------|-------|----------|
| `templated_id_basic` | `0q-nivis-zozak-QDJ-s004dev` | General resource identification |
| `infrastructure_name` | `dev-LADOZ-ZABAJ-cluster-9W7-kZ` | Infrastructure components |
| `versioned_resource_name` | `minty-v07-CngFs5K6` | Versioned deployments |
| `database_identifier` | `rural_dev_86qo8vnu` | Database naming |
| `s3_bucket_name` | `xyz-dev-eu-central-1-nomil-tiput` | S3 bucket naming |

## Usage

```bash
# Use with default values
terraform init
terraform apply

# Use the tested parameters (generates the documented outputs above)
terraform apply \
  -var="app_seed=seed-by-team-xyz" \
  -var="environment=dev" \
  -var="cluster_size=4" \
  -var="app_version=7" \
  -var="app_name=xyz" \
  -var="region=eu_central_1"

```

## Use Cases

- **Multi-environment deployments** with consistent naming patterns
- **Versioned resource management** with incremental identifiers
- **Compliance requirements** for different cloud services
- **Team collaboration** with predictable, readable identifiers
- **Infrastructure as Code** with deterministic naming

This pattern is especially useful for organizations that need consistent, predictable naming across different environments while maintaining readability and compliance with various cloud service naming requirements.
