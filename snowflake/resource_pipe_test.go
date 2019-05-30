package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccPipe(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakePipeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"snowflake_pipe", "database", "example_db"),
					resource.TestCheckResourceAttr("snowflake_pipe", "schema", "example_schema"),
					resource.TestCheckResourceAttr("snowflake_pipe", "stage", "example_stage"),
					resource.TestCheckResourceAttr("snowflake_pipe", "table", "example_table"),
					resource.TestCheckResourceAttr("snowflake_pipe", "pipe", "example_pipe"),
					resource.TestCheckResourceAttr("snowflake_pipe", "autoIngest", "true"),
				),
			},
		},
	})
}

var testSnowflakePipeConfig = `
resource "snowflake_pipe" "test" {
  database    = "example_db"
  schema      = "example_schema"
  stage       = "example_stage"
  table       = "example_table"
  pipe        = "example_pipe"
  auto_ingest = true
}
`
