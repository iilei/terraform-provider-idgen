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
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "7Kn2tT9yAR"),
				),
			},
			{
				Config: testAccNanoIDDataSourceConfigSeeded,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "578035768397"),
				),
			},
			{
				Config: testAccNanoIDDataSourceConfigGrouped,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "dl2I-NvNS-QT"),
				),
			},
			{
				Config: testAccNanoIDDataSourceConfigDashInAlphabet,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "ecdb-eafd"),
				),
			},
			{
				Config: testAccNanoIDDataSourceConfigReadableAlphabet,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "7Kn2-tT9y-AR"),
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

const testAccNanoIDDataSourceConfigDashInAlphabet = `
data "idgen_nanoid" "test" {
  length     = 9
  group_size = 4
  alphabet   = "abc-def"
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
