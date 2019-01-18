package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestRoleSnowflakeDatabase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"snowflake_role", "name", "shoprunner_terraform"),
					resource.TestCheckResourceAttr("snowflake_role", "name", "tf-test"),
				),
			},
		},
	})
}

var testSnowflakeRoleConfig = `resource "snowflake_role" "shoprunner_terraform" {
	name = "tf-test"
}`
