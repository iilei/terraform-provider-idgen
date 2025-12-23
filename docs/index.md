---
page_title: "idgen Provider"
description: |-
  The idgen provider offers flexible, human-friendly identifier generation for Terraform.
---

terraform-provider-idgen
========================

[![codecov](https://codecov.io/github/iilei/terraform-provider-idgen/graph/badge.svg?token=CZ7ZIF2FY9)](https://codecov.io/github/iilei/terraform-provider-idgen)

> [!CAUTION]
> **Not suitable for cryptographic purposes.**
>
> Do not rely on this data source when cryptographically secure random generation is required.

> [!WARNING]
> **Version 0.x Development - Breaking Changes Expected.**
>
> This provider is in initial development (0.x.x). Per [semantic versioning](https://semver.org/#spec-item-4), **breaking changes may occur in ANY release** (minor or patch) until version 1.0.0.

## Upgrading to v0.0.3

**Breaking Change:** This version migrates from `math/rand` to `math/rand/v2`, which changes the random number generation algorithm. The same seed will produce **different outputs when upgrading from v0.0.2 to v0.0.3** (though seeds remain fully deterministic within each version). Making this change now while the provider is new and the user base is small minimizes disruption before the `v1.0` release.

💥 **Seeded IDs will generate different values** compared to previous versions
  * Seeds still work deterministically within `v0.0.3` — the change only affects migration between versions
  * `proquint_canonical` is **not affected** as it strictly adheres to the canonical Proquint specification
  * Benefits: ~2x performance improvement and better statistical properties

## Motivation

The **idgen** provider offers human-friendly identifiers with knowable characteristics and a reasonable level of control over pronounceability

## Examples

Combinations of both [Proquint](https://arxiv.org/html/0901.4016) and [NanoID](https://github.com/ai/nanoid) allow for IDs that are sufficiently random-looking while remaining pronounceable.

* **ProQuint:** `babab-danol`, `kodam-kufub`, `sonop-sotof`, ...
* **Templated:** `dunov-poguv-gJqP-elfin`, `0q-zozif-zapuf-rXK-s004dev`, `snowy-v01-WazyhDQ3`, ...

## Key Features

- **Proquint and NanoID generation** ([Proquint](https://arxiv.org/html/0901.4016), [NanoID](https://github.com/ai/nanoid))
- **Templating support** to embed IDs into structured naming conventions
- **Deterministic seeding** for reproducible environments or test setups


## Quick Start

```terraform
# Generate a simple proquint identifier
data "idgen_proquint" "docs_example" {
  length = 11
  seed   = "4711"
}

output "docs_example_id" {
  value = data.idgen_proquint.docs_example.id
  # Output: babab-danol
}
```

```terraform
# Create a templated ID combining multiple types

data "idgen_templated" "docs_complex" {
  template = "{{ .proquint }}-{{ .nanoid }}-{{ .random_word }}"

  proquint = {
    seed   = "asdf"
  }

  nanoid = {
    length   = 4
    alphabet = "readable"
    seed     = "asdf"
  }

  random_word = {
    seed = "asdf"
  }
}

output "docs_complex_id" {
  value = data.idgen_templated.docs_complex.id
  # Output: dunov-poguv-gJqP-elfin
}

```

## Data Sources

- **[proquint](./data-sources/proquint)** - Pronounceable quintet identifiers
- **[proquint_canonical](./data-sources/proquint_canonical)** - IPv4/integer encoding
- **[nanoid](./data-sources/nanoid)** - URL-safe unique identifiers
- **[random_word](./data-sources/random_word)** - Dictionary-based words
- **[templated](./data-sources/templated)** - Combine multiple ID types

## Configuration

This provider requires no configuration options. Just declare it in your `required_providers` block:

```terraform
terraform {
  required_providers {
    idgen = {
      source = "iilei/idgen"
      version = "0.0.3"
      # Check https://github.com/iilei/terraform-provider-idgen/releases for the actual latest version
    }
  }
}
```

### Preflight Seed Checks

In order to be sure to yield unique ids with your specific seed on a bunch of iterations, you might want to employ
a helper script.

There is a `preflight.sh` in the [git repo](https://github.com/iilei/terraform-provider-idgen/blob/master/preflight.sh) which helps to ensure the seed doesn't yield any collusions.


```sh
./preflight.sh '#' 10000  60000
```


Output:

```
[...]
#69976 	nojus-rajop
#69977 	tovok-rajaf
#69978 	ritaz-sijal
#69979 	gopad-vuhon
#69980 	sumon-fifir
#69981 	lugug-nujoz
#69982 	puhog-vurug
#69983 	tobib-saguj
#69984 	dodij-gotip
#69985 	hosir-dozip
#69986 	puhun-rifoh
#69987 	hahuf-vivih
#69988 	vubar-mijok
#69989 	gazuf-ribut
#69990 	bimar-zudor
#69991 	hukud-kazoz
#69992 	risul-jodud
#69993 	hazoj-godub
#69994 	valoz-pitit
#69995 	gizut-kuzal
#69996 	mogun-vopur
#69997 	fufov-fahib
#69998 	dipod-pomob
#69999 	pisad-visih

# === Duplicate Analysis ===
# No duplicates found - all IDs are unique for 10000 seeds.

```
