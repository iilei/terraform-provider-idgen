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
			// Test basic unseeded generation
			{
				Config: testAccNanoIDDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_nanoid.test", "id"),
				),
			},
			// Test custom length
			{
				Config: testAccNanoIDDataSourceConfigCustomLength,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_nanoid.test", "id"),
				),
			},
			// Test seeded (deterministic) generation
			{
				Config: testAccNanoIDDataSourceConfigSeeded,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_nanoid.test", "id"),
				),
			},
			// Test with group_size
			{
				Config: testAccNanoIDDataSourceConfigGrouped,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_nanoid.test", "id"),
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
}
`

const testAccNanoIDDataSourceConfigSeeded = `
data "idgen_nanoid" "test" {
  length   = 12
  alphabet = "numeric"
  seed     = 42
}
`

const testAccNanoIDDataSourceConfigGrouped = `
data "idgen_nanoid" "test" {
  length     = 12
  group_size = 4
  alphabet   = "alphanumeric"
}
`
