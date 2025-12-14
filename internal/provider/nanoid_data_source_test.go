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
			// Test warning when alphabet contains dash and group_size is set
			{
				Config: testAccNanoIDDataSourceConfigDashInAlphabet,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_nanoid.test", "id"),
				),
			},
			// Test less_confusable alphabet
			{
				Config: testAccNanoIDDataSourceConfigLessConfusable,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_nanoid.test", "id"),
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "79ubdk868"),
				),
			},
			// Test least_confusable alphabet
			{
				Config: testAccNanoIDDataSourceConfigLeastConfusable,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_nanoid.test", "id"),
					resource.TestCheckResourceAttr("data.idgen_nanoid.test", "id", "rt4cecyqe"),
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
  seed     = "42"
}
`

const testAccNanoIDDataSourceConfigGrouped = `
data "idgen_nanoid" "test" {
  length     = 12
  group_size = 4
  alphabet   = "alphanumeric"
}
`

const testAccNanoIDDataSourceConfigDashInAlphabet = `
data "idgen_nanoid" "test" {
  length     = 12
  group_size = 4
  alphabet   = "abc-def"
}
`

const testAccNanoIDDataSourceConfigLessConfusable = `
data "idgen_nanoid" "test" {
  length   = 9
  alphabet = "less_confusable"
  seed     = "test-less-confusable"
}
`

const testAccNanoIDDataSourceConfigLeastConfusable = `
data "idgen_nanoid" "test" {
  length   = 9
  alphabet = "least_confusable"
  seed     = "test-least-confusable"
}
`
