package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNanoIDDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNanoIDDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_nanoid.test", "id"),
				),
			},
			{
				Config: testAccNanoIDDataSourceConfigCustomLength,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "CnDxXfeKNw"),
				),
			},
			{
				Config: testAccNanoIDDataSourceConfigSeeded,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "636592278400"),
				),
			},
			{
				Config: testAccNanoIDDataSourceConfigGrouped,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "MxNF-7qpU-YE"),
				),
			},
			{
				Config: testAccNanoIDDataSourceConfigReadableAlphabet,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "CnDx-XfeK-Nw"),
				),
			},
		},
	})
}

const testAccNanoIDDataSourceConfig = `
data "idgen_nanoid" "test" {}
`

const testAccNanoIDDataSourceConfigCustomLength = `
data "idgen_nanoid" "test" {
  length = 10
  seed   = "42"
}
`

const testAccNanoIDDataSourceConfigSeeded = `
data "idgen_nanoid" "test" {
  length   = 12
  alphabet = "numeric"
  seed     = "42"
}
`

const testAccNanoIDDataSourceConfigGrouped = `
data "idgen_nanoid" "test" {
  length     = 12
  group_size = 4
  alphabet   = "alphanumeric"
  seed       = "42"
}
`

const testAccNanoIDDataSourceConfigReadableAlphabet = `
data "idgen_nanoid" "test" {
  length     = 12
  group_size = 4
  alphabet   = "readable"
  seed       = "42"
}
`

func TestAccNanoIDDataSource_DashInAlphabet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNanoIDDataSourceConfigDashInAlphabet,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "dcd--fbbe"),
				),
				// Expect a warning diagnostic about dashes in custom alphabet
			},
		},
	})
}

const testAccNanoIDDataSourceConfigDashInAlphabet = `
data "idgen_nanoid" "test" {
  length     = 9
  group_size = 4
  alphabet   = "abc-def"  # Dash in alphabet + grouping = potential confusion
  seed       = "42"
}
`
