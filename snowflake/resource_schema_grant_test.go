package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSchemaGrantSnowflake(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeSchemaGrantConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "schema", "test_schema"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "database", "test_database"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "priviliges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "role", "test_role"),
				),
			},
		},
	})
}

var testSnowflakeSchemaGrantConfig = `resource "snowflake_schema_grant" "foo" {
	schema = "test_schema"
	database = "test_database"
	priviliges = ["privilege1"]
	role = "test_role"
}`
