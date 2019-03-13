package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccRoleSnowflakeDatabase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.foo", "name", "tf-test"),
				),
			},
		},
	})
}

var testSnowflakeRoleConfig = `resource "snowflake_role" "foo" {
	name = "tf-test"
}`
