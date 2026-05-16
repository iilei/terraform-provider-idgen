package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

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
							return fmt.Errorf("expected equivalent ids for 0x-prefixed and non-prefixed hex seeds, got %q and %q", idWith, idWithout)
						}

						return nil
					},
				),
			},
		},
	})
}

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
