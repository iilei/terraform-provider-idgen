# proquint-big-sha256

This example shows how to encode a SHA-256 digest using `idgen_proquint_big`.

It demonstrates that the output of `sha256()`, a `0x`-prefixed hex literal and the
same value given with `seed_format = "decimal"` all encode identically.

## Usage

```sh
terraform init
terraform apply
```

## Expected behavior

All three outputs are equal and should resolve to:

`ronum-dipuz-musad-suzop-hajad-hagiv-jukov-famog-rabag-kakog-nimil-lopis-ribib-zutod-zamab-dikot`
