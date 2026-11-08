package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccProquintBigDataSourceConfigSHA256HexMin = `
data "idgen_proquint_big" "test_min" {
  length = 95
  seed   = "0000000000000000000000000000000000000000000000000000000000000000"
}
`

const testAccProquintBigDataSourceConfigSHA256HexMax = `
data "idgen_proquint_big" "test_max" {
  length = 95
  seed   = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
}
`

const testAccProquintBigDataSourceConfigSHA256HexPrefixEquivalence = `
data "idgen_proquint_big" "with_prefix" {
	length = 95
	seed   = "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
}

data "idgen_proquint_big" "without_prefix" {
	length = 95
	seed   = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
}
`

// SHA-256("abc"), as produced by Terraform's sha256() function (lowercase hex, no prefix).
const sha256ABCProquint = "ronum-dipuz-musad-suzop-hajad-hagiv-jukov-famog-rabag-kakog-nimil-lopis-ribib-zutod-zamab-dikot"

func TestAccProquintBigDataSource_SHA256HexMin(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProquintBigDataSourceConfigSHA256HexMin,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.idgen_proquint_big.test_min",
						"seed",
						"0000000000000000000000000000000000000000000000000000000000000000",
					),
					resource.TestCheckResourceAttr(
						"data.idgen_proquint_big.test_min",
						"id",
						"babab-babab-babab-babab-babab-babab-babab-babab-babab-babab-babab-babab-babab-babab-babab-babab",
					),
				),
			},
		},
	})
}

func TestAccProquintBigDataSource_SHA256HexMax(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProquintBigDataSourceConfigSHA256HexMax,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.idgen_proquint_big.test_max",
						"seed",
						"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
					),
					resource.TestCheckResourceAttr(
						"data.idgen_proquint_big.test_max",
						"id",
						"zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz",
					),
				),
			},
		},
	})
}

func TestAccProquintBigDataSource_SHA256HexPrefixEquivalence(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProquintBigDataSourceConfigSHA256HexPrefixEquivalence,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.idgen_proquint_big.with_prefix",
						"id",
						"zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz",
					),
					resource.TestCheckResourceAttr(
						"data.idgen_proquint_big.without_prefix",
						"id",
						"zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz-zuzuz",
					),
					func(s *terraform.State) error {
						withPrefix := s.RootModule().Resources["data.idgen_proquint_big.with_prefix"]
						withoutPrefix := s.RootModule().Resources["data.idgen_proquint_big.without_prefix"]

						idWith := withPrefix.Primary.Attributes["id"]
						idWithout := withoutPrefix.Primary.Attributes["id"]

						if idWith != idWithout {
							return fmt.Errorf(
								"expected equivalent ids for 0x-prefixed and non-prefixed hex seeds, got %q and %q",
								idWith,
								idWithout,
							)
						}

						return nil
					},
				),
			},
		},
	})
}

func TestAccProquintBigDataSource_Formats(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "idgen_proquint_big" "sha" {
  length = 95
  seed   = sha256("abc")
}

data "idgen_proquint_big" "sha_decimal" {
  length      = 95
  seed_format = "decimal"
  seed        = "84342368487090800366523834928142263660104883695016514377462985829716817089965"
}

# Digits only: hex by default, so "10" is sixteen.
data "idgen_proquint_big" "digits_hex" {
  length = 11
  seed   = "10"
}

data "idgen_proquint_big" "digits_decimal" {
  length      = 11
  seed_format = "decimal"
  seed        = "10"
}

data "idgen_proquint_big" "prefixed_no_letters" {
  length = 11
  seed   = "0x10"
}

# Truncation keeps the least significant bytes (with a warning).
data "idgen_proquint_big" "truncated" {
  length = 23
  seed   = sha256("abc")
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_big.sha", "id", sha256ABCProquint),
					resource.TestCheckResourceAttr("data.idgen_proquint_big.sha_decimal", "id", sha256ABCProquint),
					resource.TestCheckResourceAttr("data.idgen_proquint_big.digits_hex", "id", "babab-babib"),
					resource.TestCheckResourceAttr("data.idgen_proquint_big.digits_decimal", "id", "babab-babap"),
					resource.TestCheckResourceAttr("data.idgen_proquint_big.prefixed_no_letters", "id", "babab-babib"),
					resource.TestCheckResourceAttr(
						"data.idgen_proquint_big.truncated",
						"id",
						"ribib-zutod-zamab-dikot",
					),
				),
			},
		},
	})
}

func TestAccProquintBigDataSource_InvalidSeeds(t *testing.T) {
	cases := map[string]struct {
		attrs string
		err   string
	}{
		"negative":       {`seed = "-5"` + "\n" + `seed_format = "decimal"`, `negative`},
		"hex as decimal": {`seed = "ff"` + "\n" + `seed_format = "decimal"`, `not a valid decimal integer`},
		"not hex":        {`seed = "xyz"`, `not a valid hex integer`},
		"unknown format": {`seed = "ff"` + "\n" + `seed_format = "base64"`, `seed_format must be`},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(
							"data \"idgen_proquint_big\" \"test\" {\n  length = 11\n  %s\n}\n",
							tc.attrs,
						),
						ExpectError: regexp.MustCompile(tc.err),
					},
				},
			})
		})
	}
}
