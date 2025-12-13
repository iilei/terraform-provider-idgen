## Template Functions

The template supports pipe-chainable string manipulation functions:

### Case Conversion
- `upper` - Convert to uppercase: `{{ .proquint | upper }}`
- `lower` - Convert to lowercase: `{{ .nanoid | lower }}`

### String Manipulation
- `replace` - Replace all occurrences: `{{ .proquint | replace "-" "_" }}`
- `substr` - Extract substring: `{{ .nanoid | substr 0 8 }}`
- `trim` - Remove leading/trailing whitespace: `{{ .value | trim }}`
- `trimPrefix` - Remove prefix: `{{ .value | trimPrefix "pre-" }}`
- `trimSuffix` - Remove suffix: `{{ .value | trimSuffix "-suf" }}`

### Repetition & Reversal
- `repeat` - Repeat string N times: `{{ "-" | repeat 3 }}`
- `reverse` - Reverse string: `{{ .random_word | reverse }}`

### Splitting & Joining
- `split` - Split into array: `{{ .value | split "-" }}`
- `join` - Join array: `{{ .array | join "_" }}`

### Examples

```hcl
# Uppercase proquint with underscores
data "idgen_templated" "example1" {
  template = "{{ .proquint | upper | replace \"-\" \"_\" }}"
  proquint = { length = 11 }
}
# Result: LUSAB_BABAD

# Complex chaining
data "idgen_templated" "example2" {
  template = "{{ .proquint | upper }}_{{ .random_word | reverse }}"
  proquint = { length = 11, seed = "test" }
  random_word = { seed = "test" }
}
# Result: TATAJ-RUBAB_tset

# Substring extraction
data "idgen_templated" "example3" {
  template = "short-{{ .nanoid | substr 0 8 | lower }}"
  nanoid = { length = 21 }
}
# Result: short-a1b2c3d4
```

