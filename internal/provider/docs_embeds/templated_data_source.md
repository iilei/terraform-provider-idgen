## Template Functions

The template supports pipe-chainable string manipulation functions:

### Case Conversion

**`upper`** - Convert to uppercase
```hcl
# Input: "sunny" | Output: "SUNNY"
template = "{{ .random_word | upper }}"
random_word = { seed = "17" }
```

**`lower`** - Convert to lowercase
```hcl
# Input: "SUNNY" | Output: "sunny"
template = "{{ .random_word | upper | lower }}"
random_word = { seed = "17" }
```

### String Manipulation

**`replace`** - Replace all occurrences
```hcl
# Input: "sunny" | Output: "sun"
template = "{{ .random_word | replace \"ny\" \"\" }}"
random_word = { seed = "17" }
```

**`substr`** - Extract substring (start, length)
```hcl
# Input: "sunny" | Output: "unn"
template = "{{ .random_word | substr 1 3 }}"
random_word = { seed = "17" }
```

**`trim`** - Remove leading/trailing whitespace
```hcl
# Input: "  sunny  " | Output: "sunny"
template = "{{ .random_word | trim }}"
random_word = { seed = "17" }
```

**`trimPrefix`** - Remove prefix
```hcl
# Input: "sunny" | Output: "nny"
template = "{{ .random_word | trimPrefix \"su\" }}"
random_word = { seed = "17" }
```

**`trimSuffix`** - Remove suffix
```hcl
# Input: "sunny" | Output: "sun"
template = "{{ .random_word | trimSuffix \"ny\" }}"
random_word = { seed = "17" }
```

### Repetition & Reversal

**`repeat`** - Repeat string N times
```hcl
# Input: "sunny" | Output: "sunnysunnysunny"
template = "{{ .random_word | repeat 3 }}"
random_word = { seed = "17" }
```

**`reverse`** - Reverse string
```hcl
# Input: "sunny" | Output: "ynnus"
template = "{{ .random_word | reverse }}"
random_word = { seed = "17" }
```

### More Examples

```hcl
# yields: 0q-LUSAB_BABAD
data "idgen_templated" "example1" {
  template = "0q-{{ .proquint_canonical | upper | replace \"-\" \"_\" }}"
  proquint_canonical = { seed = "127.0.0.1" }
}

# yields: snowy-vibub-vamiz
data "idgen_templated" "example2" {
  template = "{{ .random_word }}-{{ .proquint }}"
  random_word = { seed = "a5e57e8a-9a7c-4efd-9fdd-0fcdc7630e3a" }
  proquint = { seed = "a5e57e8a-9a7c-4efd-9fdd-0fcdc7630e3a" }
}

# yields: bCi-pPV
data "idgen_templated" "example3" {
  template = "{{ .nanoid }}"
  nanoid = { length = 21, seed = "72da0233-3b03-4410-854f-3b96e868e15a", alphabet = "readable", length = 7, group_size = 3 }
}


```

