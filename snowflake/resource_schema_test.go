package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSchema(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeSchemaConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"snowflake_schema", "database", "example_db"),
					resource.TestCheckResourceAttr("snowflake_schema", "schema", "example_schema"),
				),
			},
		},
	})
}

var testSnowflakeSchemaConfig = `
resource "snowflake_schema" "test" {
  database = "example_db"
  schema   = "example_schema"
}
`
