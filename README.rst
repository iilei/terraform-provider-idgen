terraform-provider-idgen
=======================

.. image:: https://codecov.io/github/iilei/terraform-provider-idgen/graph/badge.svg?token=CZ7ZIF2FY9 
 :target: https://codecov.io/github/iilei/terraform-provider-idgen
 
.. warning::
   ⚠️ **Version 0.x Development - Breaking Changes Expected**

   This provider is in initial development (0.x.x). Per `semantic versioning <https://semver.org/#spec-item-4>`_, **breaking changes may occur in ANY release** (minor or patch) until version 1.0.0.

.. warning::
  ⚠️ **Not suitable for cryptographic purposes**

  Do not rely on this data source when cryptographically secure random generation is required.


The **idgen** provider offers human-friendly identifiers with knowable characteristics and a reasonable level of control over pronounceability.

Combinations of both `Proquint <https://arxiv.org/html/0901.4016>`_ and `NanoID <https://github.com/ai/nanoid>`_ allow for IDs that are sufficiently random-looking whilst at the same time being pronounceable.


Key Features
------------

- **Proquint and NanoID generation** (`Proquint <https://arxiv.org/html/0901.4016>`_, `NanoID <https://github.com/ai/nanoid>`_)
- **Templating support** to embed IDs into structured naming conventions
- **Deterministic seeding** for reproducible environments or test setups

Example
-------

Basic ID Generation
~~~~~~~~~~~~~~~~~~~

.. code-block:: hcl

   # Generate a NanoID with a total length of 12, grouped every 4 characters
   data "idgen_nanoid" "example" {
     length     = 12
     group_size = 4
     alphabet   = "readable"
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
If desired, this can be easily added as follows:

.. code-block:: hcl

  # Generate a single templated ID combining Proquint, NanoID and local variables.
  # Seeds result in fully deterministic IDs.
  locals {
    # the seed is going to produce deterministic IDs for both nanoid and proquint
    seed = "app-specific-seed"

    # other variables
    size     = 4
    size_fmt = format("%03d", local.size)  # "004"
    stage    = "dev"
  }

  data "idgen_templated" "example" {
    template = "0q-{{ .proquint }}-{{ .nanoid }}-s${local.size_fmt}-${local.stage}"
    nanoid   = { length = 3, seed = "#${local.size}_${local.seed}", alphabet = "readable" }
    proquint = { length = 11, seed = "#${local.size}_${local.seed}" }
  }

  output "my_templated_id" {
    value = data.idgen_templated.example.id
  }

  # yields: "0q-zozif-zapuf-rXK-s004-dev"

Random Word Generation
~~~~~~~~~~~~~~~~~~~~~~

Picks a word based on ``seed`` and ``wordlist``. The default word list is intentionally kept short and will not be extended further.

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
       random        = data.idgen_random_word.random.id
       deterministic = data.idgen_random_word.deterministic.id
       custom        = data.idgen_random_word.custom.id
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
        seed   = "my-app-42"  # => deterministic but random-looking
      }

      data "idgen_nanoid" "user_id" {
        length = 12
        seed   = "user-alice"  # => deterministic NanoID
      }

**Without Seed (Random)**
   Omitting ``seed`` generates more random numbers, encoded as ProQuint strings.

   .. code-block:: hcl

      data "idgen_proquint" "random" {
        # No seed => some random proquint within the range from *babab-babab* (`0`) to *zuzuz-zuzuz* (`4_294_967_295`)
        length = 11
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

  TF_VAR_seed_prefix="" ./test-provider.sh

Example output:

.. code-block:: console

  Seed                           | NanoID          | Proquint
  ------------------------------------------------------------------------
  asdf                           | M0o-t2I         | dunov-poguv
  asdf-1                         | Q5y-LKz         | nizik-hojiz
  [...]
  asdf-11                        | 1UP-JVm         | bobud-dahip
  asdf-12                        | A4d-cFi         | vatut-kuvag
