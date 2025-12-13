terraform-provider-idgen
=======================

.. warning::
   ⚠️ **Version 0.x Development - Breaking Changes Expected**

   This provider is in initial development (0.x.x). Per `semantic versioning <https://semver.org/#spec-item-4>`_, **breaking changes may occur in ANY release** (minor or patch) until version 1.0.0.

.. caution::
   ⚠️ **Security Notice: Seeded IDs**

   When using a **seed** parameter, IDs become **deterministic and predictable**. Seeded IDs should **never** be used for security tokens, passwords, secrets, session IDs, or any cryptographic purpose. Use seeding only for reproducible naming in test environments or infrastructure patterns.


The **idgen** provider offers flexible, human-friendly identifier generation for Terraform.
It supports multiple ID formats including **Proquint** and **NanoID**, with optional **templating**, **controlled entropy**, and **seed-based determinism**.
These IDs are read-only utilities for use within Terraform configurations, making them ideal for predictable, human-readable identifiers without managing lifecycle resources.

Key Features
------------

- **Proquint and NanoID generation** (`Proquint <https://arxiv.org/html/0901.4016>`_, `NanoID <https://github.com/ai/nanoid>`_)
- **Configurable entropy** for predictable or high-randomness IDs
- **Templating support** to embed IDs into structured naming conventions
- **Deterministic seeding** for reproducible environments or test setups
- **Terraform-native usage as data sources** — no resource lifecycle management required

Example
-------

Basic ID Generation
~~~~~~~~~~~~~~~~~~~

.. code-block:: hcl

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
   }

   output "my_ids" {
     value = {
       nanoid   = data.idgen_nanoid.example.id
       proquint = data.idgen_proquint.example.id
     }
   }

Templated IDs with Parametrization
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The `Proquint specification <https://arxiv.org/html/0901.4016#_conclusion_and_specification>`_ suggests using the optional magic number prefix ``0q-`` before a sequence of proquints for clarity.

.. code-block:: hcl

   # Generate a single templated ID combining NanoID and Proquint, fully deterministic.
   # yields the following exact string: 0q-TATAJ-RUBAB.Yuc.DZH.5iW

   data "idgen_templated" "example" {
     template = "0q-{{ .proquint | upper }}.{{ .nanoid | replace '-' '.' }}"

    # the seed is going to produce deterministic IDs for both nanoid and proquint
     proquint = {
       length = 9
       seed   = "xyz-12"  # => tataj-rubab
       group_size = 3
     }

     nanoid = {
       length = 9
       seed   = "xyz-12"  # => Yuc-DZH-5iW
       group_size = 3
       alphabet = "readable"
     }
   }

   output "my_templated_id" {
     value = data.idgen_templated.example.id
   }

Random Word Generation
~~~~~~~~~~~~~~~~~~~~~~

Generate human-friendly identifiers using random words from a curated five-letter word list:

.. code-block:: hcl

   # Generate a random word (changes on each apply)
   data "idgen_random_word" "random" {}

   # Generate a deterministic word using a seed
   data "idgen_random_word" "deterministic" {
     seed = "some-seed"  # Always produces the same word
   }

   # Use a custom word list
   data "idgen_random_word" "custom" {
     seed     = "5"  # Produces "red" from the custom list
     wordlist = "red,blue,green,yellow,purple"
   }

   output "identifiers" {
     value = {
       random       = data.idgen_random_word.random.id
       deterministic = data.idgen_random_word.deterministic.id
       custom       = data.idgen_random_word.custom.id
     }
   }

.. note::
   Sequential numeric seeds (``"0"``, ``"1"``, ``"2"``) produce words in alphabetical order. For varied distribution, use non-numeric seeds (``"project-1"``, ``"env-prod"``) which are hashed for randomized selection.

Alphabet Presets
----------------

NanoID supports configurable alphabets. You can either use **named presets** for convenience or provide a **custom string** of allowed characters.

.. list-table::
   :header-rows: 1
   :widths: 20 80

   * - Preset Name
     - Description
   * - ``alphanumeric``
     - Uppercase + lowercase letters and digits (``a-zA-Z0-9``)
   * - ``numeric``
     - Digits only (``0-9``)
   * - ``readable``
     - Avoids visually confusing characters (e.g., ``0/O``, ``1/l``)


Seed Parameter Behavior
~~~~~~~~~~~~~~~~~~~~~~~

The ``seed`` parameter provides **deterministic ID generation** with smart behavior based on input format:

**IPv4 Addresses (Direct Encoding)**
   IPv4 addresses are directly encoded as proquints, following the canonical specification:

   .. code-block:: hcl

      data "idgen_proquint" "localhost" {
        seed = "127.0.0.1"  # => lusab-babad
      }

**Integers in uint32 Range (Direct Encoding)**
   Integers from 0 to 4,294,967,295 are directly encoded:

   .. code-block:: hcl

      data "idgen_proquint" "from_number" {
        seed = "2130706433"  # => lusab-babad (same as 127.0.0.1)
      }

**Text Strings (Seeded Random Generation)**
   Any string (or large integer) is hashed to create a seed for random generation:

   .. code-block:: hcl

      data "idgen_proquint" "app_id" {
        length = 17
        seed   = "my-app-42"  # => deterministic but random-looking ID
      }

      data "idgen_nanoid" "user_id" {
        length = 12
        seed   = "user-alice"  # => deterministic NanoID
      }

**Without Seed (Cryptographically Random)**
   Omitting ``seed`` generates cryptographically secure random IDs:

   .. code-block:: hcl

      data "idgen_proquint" "random" {
        # No seed => different ID on each apply
      }

.. note::
   proquint seeds are treated as numbers or IPv4 addresses when possible for canonical behavior. For that reason, if high entropy is desired, add a non-numeric part to the seed string to force random generation.

Notes
~~~~~

- ``length`` controls the **total number of characters**
- ``group_size`` defines how many characters are per group split by dash (``-``)
- ``alphabet`` supports **named presets** for ease of use, or users can provide a custom string.
- The ``idgen_templated`` data source allows **parametrized combination** of multiple base IDs, with optional inline transformations (``upper``, ``lower``, etc.)
- Terraform-native string interpolation can still be used for additional customization if needed.


Local Development
-----------------

Generate a few IDs with different seeds to see how the provider behaves:

.. code-block:: bash

    ./test-provider.sh

Example output:

.. code-block:: console

  Testing ID generation with different seeds
  ==========================================
  Seed                           | NanoID          | Proquint
  ------------------------------------------------------------------------
  asdf                           | M0o-t2I-uP3     | dunov-poguv
  asdf-1                         | Q5y-LKz-mwJ     | nizik-hojiz
  asdf-2                         | kWB-v3C-ZuP     | gufat-horub
  asdf-3                         | L3n-wgr-rZb     | bozag-jibad
  asdf-4                         | HPP-8Tf-ED1     | gapuk-ginop
  asdf-5                         | hRa-puB-Sq0     | makir-zabit
  asdf-6                         | RUa-kAH-Rce     | fijif-gakoj
  asdf-7                         | TA3-HkY-YRy     | kodam-kufub
  asdf-8                         | qNn-sE7-k5V     | nufuv-hosos
  asdf-9                         | Pyi-5lT-OwP     | junoh-bizah
  asdf-10                        | 1d8-7Kd-FYD     | sonop-sotof
  asdf-11                        | 1UP-JVm-Eyc     | bobud-dahip
  asdf-12                        | A4d-cFi-2pw     | vatut-kuvag

  ==========================================
  Try different seed prefixes:
    TF_VAR_seed_prefix=myapp ./test-provider.sh
    TF_VAR_seed_prefix=127.0.0.1 ./test-provider.sh
    TF_VAR_seed_prefix=$(date +"x%s%N") ./test-provider.sh
  ==========================================


..
   internal notes:
   named alphabet presets for nanoid and proquint:
   https://github.com/matoous/go-nanoid/blob/main/gonanoid.go#L9-L39
