package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTemplatedDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccTemplatedDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_templated.test", "id"),
					resource.TestCheckResourceAttr("data.idgen_templated.test", "id", "fixed-templated-thing-abc-123"),
				),
			},
		},
	})
}

const testAccTemplatedDataSourceConfig = `
data "idgen_templated" "test" {}
`
