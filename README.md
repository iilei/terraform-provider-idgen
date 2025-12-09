# terraform-provider-idgen

> ⚠️ **Version 0.x Development - Breaking Changes Expected**
>
> This provider is in initial development (0.x.x). Per [semantic versioning](https://semver.org/#spec-item-4), **breaking changes may occur in ANY release** (minor or patch) until version 1.0.0.


The **idgen** provider offers flexible, human-friendly identifier generation for Terraform.
It supports multiple ID formats including **Proquint** and **NanoID**, with optional **templating**, **controlled entropy**, and **seed-based determinism**.
These IDs are read-only utilities for use within Terraform configurations, making them ideal for predictable, human-readable identifiers without managing lifecycle resources.

## Key Features

- **Proquint and NanoID generation**
- **Configurable entropy** for predictable or high-randomness IDs
- **Templating support** to embed IDs into structured naming conventions
- **Deterministic seeding** for reproducible environments or test setups
- **Terraform-native usage as data sources** — no resource lifecycle management required

## Example

### Basic ID Generation

```hcl
# Generate a NanoID with a total length of 12, grouped every 4 characters
data "idgen_nanoid" "example" {
  length     = 12
  group_size = 4
  alphabet   = "alphanumeric" # preset: a-zA-Z0-9
}

# Generate a Proquint ID with a total length of 12 (entropy calculated internally), grouped every 4 characters
data "idgen_proquint" "example" {
  length     = 12
  group_size = 4
  alphabet   = "standard" # preset for Proquint
}

output "my_ids" {
  value = {
    nanoid   = data.idgen_nanoid.example.id
    proquint = data.idgen_proquint.example.id
  }
}
```

### Templated IDs with Parametrization

```hcl
# Generate a single templated ID combining NanoID and Proquint, fully deterministic
data "idgen_templated" "example" {
  template = "thing-{{ .proquint | upper }}.{{ .nanoid | replace '-' '.' }}"

  nanoid = {
    length = 6
    seed   = 42
    group_size = 3
    alphabet = "numeric"
  }

  proquint = {
    length = 9
    seed   = 42
    group_size = 3
    alphabet = "standard"
  }
}

output "my_templated_id" {
  value = data.idgen_templated.example.id
}

```

## Alphabet Presets

Both NanoID and Proquint support configurable alphabets. You can either use **named presets** for convenience or provide a **custom string** of allowed characters.

| Preset Name      | Description |
|------------------|-------------|
| `alphanumeric`   | Uppercase + lowercase letters and digits (`a-zA-Z0-9`) |
| `numeric`        | Digits only (`0-9`) |
| `readable`       | Avoids visually confusing characters (e.g., `0/O`, `1/l`) |



### Notes

- `length` controls the **total number of characters**
- `group_size` defines how many characters are per group split by dash (`-`)
- `alphabet` supports **named presets** for ease of use, or users can provide a custom string.
- The `idgen_templated` data source allows **parametrized combination** of multiple base IDs, with optional inline transformations (`upper`, `lower`, etc.)

- Terraform-native string interpolation can still be used for additional customization if needed.
