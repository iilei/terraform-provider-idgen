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
			// Test basic template with proquint and nanoid
			{
				Config: testAccTemplatedDataSourceConfigBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.idgen_templated.test", "id", "babab-babad-babab-babab-aysZ"),
					// Just verify it's set and has expected structure
				),
			},
			// Test template with random_word
			{
				Config: testAccTemplatedDataSourceConfigWithWord,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_templated.test", "id"),
					resource.TestCheckResourceAttr("data.idgen_templated.test", "id", "apple-tataj-rubab"),
				),
			},
			// Test template with all ID types combined
			{
				Config: testAccTemplatedDataSourceConfigAllTypes,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_templated.test", "id"),
					resource.TestCheckResourceAttr("data.idgen_templated.test", "id", "lusab-babad.tataj-rubab.ays-ZFz-mgE-es.apple"),
				),
			},
			// Test template functions with piping
			{
				Config: testAccTemplatedDataSourceConfigWithFunctions,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.idgen_templated.test", "id"),
					resource.TestCheckResourceAttr("data.idgen_templated.test", "id", "TATAJ_RUBAB_:apfel-:apfel-"),
				),
			},
		},
	})
}

const testAccTemplatedDataSourceConfigBasic = `
data "idgen_templated" "test" {
  template = "{{ .proquint_canonical }}-{{ .nanoid }}"

  proquint_canonical = {
    seed = "4294967296"
  }


  nanoid = {
    length = 4
    seed   = "xyz-12"
    alphabet = "readable"
  }
}
`

const testAccTemplatedDataSourceConfigWithWord = `
data "idgen_templated" "test" {
  template = "{{ .random_word }}-{{ .proquint }}"

  random_word = {
    seed     = "0"
    wordlist = "apple,banana,cherry"
  }

  proquint = {
    length = 11
    seed   = "xyz-12"
    group_size = 5
  }
}
`

const testAccTemplatedDataSourceConfigAllTypes = `
data "idgen_templated" "test" {
  template = "{{ .proquint_canonical }}.{{ .proquint }}.{{ .nanoid }}.{{ .random_word }}"

  proquint_canonical = {
    seed = "127.0.0.1"
  }

  proquint = {
    length = 11
    seed   = "xyz-12"
  }

  nanoid = {
    length = 14
    seed   = "xyz-12"
    alphabet = "readable"
    group_size = 3
  }

  random_word = {
    seed     = "0"
    wordlist = "apple,banana,cherry"
  }
}
`

const testAccTemplatedDataSourceConfigWithFunctions = `
data "idgen_templated" "test" {
  template = "{{ .proquint | upper | replace \"-\" \"_\" }}_{{ .random_word | reverse | replace \"elppa\" \"apfel\" | append \"-\" | prepend \":\" | repeat 2 }}"

  proquint = {
    length = 11
    seed   = "xyz-12"
  }

  random_word = {
    seed     = "0"
    wordlist = "apple,banana,cherry"
  }
}
`
