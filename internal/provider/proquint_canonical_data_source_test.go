package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProquintCanonicalDataSource_IPv4(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProquintCanonicalDataSourceConfig_IPv4("127.0.0.1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "127.0.0.1"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "lusab-babad"),
				),
			},
		},
	})
}

func TestAccProquintCanonicalDataSource_IPv4Examples(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProquintCanonicalDataSourceConfig_Multiple(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Examples from the original proquint paper
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.localhost", "id", "lusab-babad"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.google", "id", "gutih-tugad"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.max", "id", "zuzuz-zuzuz"),
				),
			},
		},
	})
}

func TestAccProquintCanonicalDataSource_Uint32(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProquintCanonicalDataSourceConfig_Uint32("2130706433"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "2130706433"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "lusab-babad"),
				),
			},
			{
				Config: testAccProquintCanonicalDataSourceConfig_Uint32("0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "0"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "babab-babab"),
				),
			},
			{
				Config: testAccProquintCanonicalDataSourceConfig_Uint32("4294967295"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "4294967295"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "zuzuz-zuzuz"),
				),
			},
		},
	})
}

func TestAccProquintCanonicalDataSource_Uint64(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProquintCanonicalDataSourceConfig_Uint32("4294967296"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "4294967296"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "babab-babad-babab-babab"),
				),
			},
			{
				Config: testAccProquintCanonicalDataSourceConfig_Uint32("9223372036854775807"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "9223372036854775807"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "luzuz-zuzuz-zuzuz-zuzuz"),
				),
			},
			{
				Config: testAccProquintCanonicalDataSourceConfig_Uint32("18446744073709551615"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "18446744073709551615"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "zuzuz-zuzuz-zuzuz-zuzuz"),
				),
			},
		},
	})
}

func TestAccProquintCanonicalDataSource_Hexadecimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProquintCanonicalDataSourceConfig_Uint32("0x7f000001"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "0x7f000001"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "lusab-babad"),
				),
			},
			{
				Config: testAccProquintCanonicalDataSourceConfig_Uint32("7f000001"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "7f000001"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "lusab-babad"),
				),
			},
			{
				Config: testAccProquintCanonicalDataSourceConfig_Uint32("0xFFFFFFFF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "0xFFFFFFFF"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "zuzuz-zuzuz"),
				),
			},
			{
				Config: testAccProquintCanonicalDataSourceConfig_Uint32("0x7fffffffffffffff"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "seed", "0x7fffffffffffffff"),
					resource.TestCheckResourceAttr("data.idgen_proquint_canonical.test", "id", "luzuz-zuzuz-zuzuz-zuzuz"),
				),
			},
		},
	})
}

func TestAccProquintCanonicalDataSource_InvalidValue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccProquintCanonicalDataSourceConfig_IPv4("not-an-ip"),
				ExpectError: regexp.MustCompile("Invalid seed for canonical encoding"),
			},
			{
				Config:      testAccProquintCanonicalDataSourceConfig_Uint32("99999999999999999999"),
				ExpectError: regexp.MustCompile("Invalid seed for canonical encoding"),
			},
			{
				Config:      testAccProquintCanonicalDataSourceConfig_IPv4("2001:db8::1"),
				ExpectError: regexp.MustCompile("Invalid seed for canonical encoding"),
			},
		},
	})
}

func testAccProquintCanonicalDataSourceConfig_IPv4(value string) string {
	return `
data "idgen_proquint_canonical" "test" {
  seed = "` + value + `"
}
`
}

func testAccProquintCanonicalDataSourceConfig_Uint32(value string) string {
	return `
data "idgen_proquint_canonical" "test" {
  seed = "` + value + `"
}
`
}

func testAccProquintCanonicalDataSourceConfig_Multiple() string {
	return `
data "idgen_proquint_canonical" "localhost" {
  seed = "127.0.0.1"
}

data "idgen_proquint_canonical" "google" {
  seed = "63.84.220.193"
}

data "idgen_proquint_canonical" "max" {
  seed = "255.255.255.255"
}
`
}
