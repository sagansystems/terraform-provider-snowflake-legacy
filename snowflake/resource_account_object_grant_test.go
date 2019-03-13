package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAccountObjectGrantSnowflake(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeAccountObjectGrantConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_object_grant.foo", "object_type", "test_type"),
					resource.TestCheckResourceAttr("snowflake_account_object_grant.foo", "object_name", "test_name"),
					resource.TestCheckResourceAttr("snowflake_account_object_grant.foo", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_account_object_grant.foo", "grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_account_object_grant.foo", "role", "test_role"),
				),
			},
		},
	})
}

var testSnowflakeAccountObjectGrantConfig = `resource "snowflake_account_object_grant" "foo" {
	object_type = "test_type"
	object_name = "test_name"
	privileges = ["privilege1"]
	role = "test_role"
}`
