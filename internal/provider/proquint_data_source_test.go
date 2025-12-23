package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccProquintDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test basic unseeded generation
			{
				Config: testAccProquintDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_proquint.test", "id"),
				),
			},
			{
				Config: testAccProquintDataSourceConfigAsDocumented,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.localhost", "id", "lusab-babad"),
				),
			},
			// Test seeded (deterministic) generation
			{
				Config: testAccProquintDataSourceConfigSeeded,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.test", "id", "lufuh-fumod-tagan"),
				),
			},
			// Test with group_size (unseeded should produce different results)
			{
				Config: testAccProquintDataSourceConfigGrouped,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_proquint.test_a", "id"),
					resource.TestCheckResourceAttrSet("data.idgen_proquint.test_b", "id"),
					func(s *terraform.State) error {
						rsA := s.RootModule().Resources["data.idgen_proquint.test_a"]
						rsB := s.RootModule().Resources["data.idgen_proquint.test_b"]
						if rsA.Primary.Attributes["id"] == rsB.Primary.Attributes["id"] {
							return fmt.Errorf("unseeded proquints should differ, but both are: %s", rsA.Primary.Attributes["id"])
						}
						return nil
					},
				),
			},
			// Test deterministic seeded generation (same seed = same output)
			{
				Config: testAccProquintDataSourceConfigSeededDeterministic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.test1", "id", "lufuh-fumod"),
					resource.TestCheckResourceAttr("data.idgen_proquint.test2", "id", "lufuh-fumod"),
				),
			},
			// Test string seed (non-numeric)
			{
				Config: testAccProquintDataSourceConfigStringSeed,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.test", "id", "rofuz-ropot"),
				),
			},
			// Test minimum length
			{
				Config: testAccProquintDataSourceConfigMinLength,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.test", "id", "ludub"),
				),
			},
			// Test seeded with grouping
			{
				Config: testAccProquintDataSourceConfigSeededGrouped,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.test", "id", "jol-abl-iba-rlo-fas"),
				),
			},
			// Test group_size = 0 (no grouping)
			{
				Config: testAccProquintDataSourceConfigZeroGroupSize,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.test", "id", "kasoj-vizah"),
				),
			},
			// Test small group_size
			{
				Config: testAccProquintDataSourceConfigSmallGroupSize,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.test", "id", "bi-zi-mv-ij-iv"),
				),
			},
			// Test invalid length (should fail)
			{
				Config:      testAccProquintDataSourceConfigInvalidLength,
				ExpectError: regexp.MustCompile("Length exceeds maximum allowed value"),
			},
			// Test various IP addresses (direct encoding)
			{
				Config: testAccProquintDataSourceConfigIPAddresses,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.ip1", "id", "lusab-babad"),
					resource.TestCheckResourceAttr("data.idgen_proquint.ip2", "id", "zusab-babab"),
					resource.TestCheckResourceAttr("data.idgen_proquint.ip3", "id", "zusab-zusad"),
				),
			},
			// Test various IP addresses with truncation
			{
				Config: testAccProquintDataSourceConfigIPAddressesTrunc,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.ip1_trunc", "id", "babad"),
					resource.TestCheckResourceAttr("data.idgen_proquint.ip2_trunc", "id", "babab"),
					resource.TestCheckResourceAttr("data.idgen_proquint.ip3_trunc", "id", "zusad"),
				),
			},
			// Test various IP addresses with padding
			{
				Config: testAccProquintDataSourceConfigIPAddressesPad,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.ip1_pad", "id", "babab-lusab-babad"),
					resource.TestCheckResourceAttr("data.idgen_proquint.ip2_pad", "id", "babab-zusab-babab"),
					resource.TestCheckResourceAttr("data.idgen_proquint.ip3_pad", "id", "babab-zusab-zusad"),
				),
			},
			// Test decimal uint32 (directly encoded like IP - canonical proquint behavior)
			{
				Config: testAccProquintDataSourceConfigDecimalUint,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.decimal", "id", "lusab-babad"),
				),
			},
			// Test that IP and its decimal equivalent produce the same result
			{
				Config: testAccProquintDataSourceConfigIPEquivalence,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint.from_ip", "id", "lusab-babad"),
				),
			},
		},
	})
}

const testAccProquintDataSourceConfig = `
data "idgen_proquint" "test" {
  length = 11
}
`

const testAccProquintDataSourceConfigAsDocumented = `
data "idgen_proquint" "localhost" {
  length = 11
  seed   = "127.0.0.1"
}
`

const testAccProquintDataSourceConfigSeeded = `
data "idgen_proquint" "test" {
  length = 17
  seed   = "seed-42"
}
`

const testAccProquintDataSourceConfigGrouped = `
data "idgen_proquint" "test_a" {
  length     = 17
  group_size = 5
}
data "idgen_proquint" "test_b" {
  length     = 17
  group_size = 5
}
`

const testAccProquintDataSourceConfigSeededDeterministic = `
data "idgen_proquint" "test1" {
  length = 11
  seed   = "seed-42"
}

data "idgen_proquint" "test2" {
  length = 11
  seed   = "seed-42"
}
`

const testAccProquintDataSourceConfigStringSeed = `
data "idgen_proquint" "test" {
  length = 11
  seed   = "production-env"
}
`

const testAccProquintDataSourceConfigMinLength = `
data "idgen_proquint" "test" {
  length = 5
  seed   = "*"
}
`

const testAccProquintDataSourceConfigSeededGrouped = `
data "idgen_proquint" "test" {
  length     = 17
  seed       = "grouped"
  group_size = 3
}
`

const testAccProquintDataSourceConfigZeroGroupSize = `
data "idgen_proquint" "test" {
  length     = 11
  seed       = "zero"
  group_size = 0
}
`

const testAccProquintDataSourceConfigSmallGroupSize = `
data "idgen_proquint" "test" {
  length     = 11
  seed       = "small"
  group_size = 2
}
`

const testAccProquintDataSourceConfigInvalidLength = `
data "idgen_proquint" "test" {
  length = 2000
}
`

const testAccProquintDataSourceConfigIPAddresses = `
data "idgen_proquint" "ip1" {
  length = 11
  seed   = "127.0.0.1"
}

data "idgen_proquint" "ip2" {
  length = 11
  seed   = "255.0.0.0"
}

data "idgen_proquint" "ip3" {
  length = 11
  seed   = "255.0.255.1"
}
`

const testAccProquintDataSourceConfigIPAddressesTrunc = `
data "idgen_proquint" "ip1_trunc" {
  length = 5
  seed   = "127.0.0.1"
}

data "idgen_proquint" "ip2_trunc" {
  length = 5
  seed   = "255.0.0.0"
}

data "idgen_proquint" "ip3_trunc" {
  length = 5
  seed   = "255.0.255.1"
}
`

const testAccProquintDataSourceConfigIPAddressesPad = `
data "idgen_proquint" "ip1_pad" {
  length = 17
  seed   = "127.0.0.1"
}

data "idgen_proquint" "ip2_pad" {
  length = 17
  seed   = "255.0.0.0"
}

data "idgen_proquint" "ip3_pad" {
  length = 17
  seed   = "255.0.255.1"
}
`

const testAccProquintDataSourceConfigDecimalUint = `
data "idgen_proquint" "decimal" {
  length = 11
  seed   = "2130706433"
}
`

const testAccProquintDataSourceConfigIPEquivalence = `
# IPv4 addresses are directly encoded, demonstrating the canonical proquint behavior
data "idgen_proquint" "from_ip" {
  length = 11
  seed   = "127.0.0.1"
}
`
