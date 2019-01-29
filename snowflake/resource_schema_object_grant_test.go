package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestSchemaObjectGrantSnowflake(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeSchemaObjectGrantConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "object_type", "TABLE"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "object_name", "SAMPLE_TABLE"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "database", "MASTER"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "schema", "SAMPLE_SCHEMA"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "future", "false"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "priviliges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.foo", "role", "test_role"),
				),
			},
		},
	})
}

var testSnowflakeSchemaObjectGrantConfig = `resource "snowflake_schema_object_grant" "foo" {
	object_type = "TABLE"
	object_name = "SAMPLE_TABLE"
	database = "MASTER"
	schema = "SAMPLE_SCHEMA""
	priviliges = ["privilege1"]
	role = "SAMPLE_ROLE"
}`
