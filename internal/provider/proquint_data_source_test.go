package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
			// Test seeded (deterministic) generation
			{
				Config: testAccProquintDataSourceConfigSeeded,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_proquint.test", "id"),
				),
			},
		},
	})
}

const testAccProquintDataSourceConfig = `
data "idgen_proquint" "test" {}
`

const testAccProquintDataSourceConfigSeeded = `
data "idgen_proquint" "test" {
  length = 17
  seed   = 42
}
`
