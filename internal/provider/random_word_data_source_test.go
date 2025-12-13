package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRandomWordDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test basic unseeded generation (random word)
			{
				Config: testAccRandomWordDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_random_word.test", "id"),
				),
			},
			// Test seeded (deterministic) generation with numeric seed
			{
				Config: testAccRandomWordDataSourceConfigSeededNumeric,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_random_word.test", "id"),
					resource.TestCheckResourceAttr("data.idgen_random_word.test", "seed", "0"),
				),
			},
			// Test seeded generation with text seed (hashed)
			{
				Config: testAccRandomWordDataSourceConfigSeededText,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_random_word.test", "id"),
					resource.TestCheckResourceAttr("data.idgen_random_word.test", "seed", "my-project-seed"),
				),
			},
			// Test custom wordlist
			{
				Config: testAccRandomWordDataSourceConfigCustomWordlist,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_random_word.test", "id", "banana"),
					resource.TestCheckResourceAttr("data.idgen_random_word.test", "seed", "1"),
					resource.TestCheckResourceAttr("data.idgen_random_word.test", "wordlist", "apple,banana,cherry,date"),
				),
			},
			// Test custom wordlist with text seed
			{
				Config: testAccRandomWordDataSourceConfigCustomWordlistWithTextSeed,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_random_word.test", "id", "elderberry"),
					resource.TestCheckResourceAttr("data.idgen_random_word.test", "seed", "my-project-seed-1"),
					resource.TestCheckResourceAttr("data.idgen_random_word.test", "wordlist", "apple,banana,cherry,date,elderberry"),
				),
			},
		},
	})
}

const testAccRandomWordDataSourceConfig = `
data "idgen_random_word" "test" {}
`

const testAccRandomWordDataSourceConfigSeededNumeric = `
data "idgen_random_word" "test" {
  seed = "0"
}
`

const testAccRandomWordDataSourceConfigSeededText = `
data "idgen_random_word" "test" {
  seed = "my-project-seed"
}
`

const testAccRandomWordDataSourceConfigCustomWordlist = `
data "idgen_random_word" "test" {
  seed     = "1"
  wordlist = "apple,banana,cherry,date"
}
`

const testAccRandomWordDataSourceConfigCustomWordlistWithTextSeed = `
data "idgen_random_word" "test" {
  seed     = "my-project-seed-1"
  wordlist = "apple,banana,cherry,date,elderberry"
}
`
