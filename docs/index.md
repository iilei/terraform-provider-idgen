---
page_title: "idgen Provider"
description: |-
  The idgen provider offers flexible, human-friendly identifier generation for Terraform.
---

terraform-provider-idgen
=======================

[![codecov](https://codecov.io/github/iilei/terraform-provider-idgen/graph/badge.svg?token=CZ7ZIF2FY9)](https://codecov.io/github/iilei/terraform-provider-idgen)

> [!CAUTION]
> **Not suitable for cryptographic purposes.**
>
> Do not rely on this data source when cryptographically secure random generation is required.

> [!WARNING]
> **Version 0.x Development - Breaking Changes Expected.**
>
> This provider is in initial development (0.x.x). Per [semantic versioning](https://semver.org/#spec-item-4), **breaking changes may occur in ANY release** (minor or patch) until version 1.0.0.

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
      version = "0.0.2"
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
  #69979 	nilij-muzal
  #69980 	dunis-nihag
  #69981 	vakad-luzaz
  #69982 	mibum-kalas
  #69983 	jitov-dozan
  #69984 	zosoh-zugaj
  #69985 	sorap-niron
  #69986 	lanut-fizoh
  #69987 	hahug-haror
  #69988 	vilon-dajos
  #69989 	pizaz-nosan
  #69990 	rikik-savom
  #69991 	lavoj-mokal
  #69992 	gumup-vogus
  #69993 	luzah-zukiv
  #69994 	togid-vufop
  #69995 	bolom-zapin
  #69996 	vihuh-latut
  #69997 	datuz-jivah
  #69998 	kumam-bodod
  #69999 	nubof-julib
  # === Duplicate Analysis ===
  # No duplicates found - all IDs are unique for 10000 seeds.

```
