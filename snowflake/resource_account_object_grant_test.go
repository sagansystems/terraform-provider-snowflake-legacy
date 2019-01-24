package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccountObjectGrantSnowflakeDatabase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"snowflake_account_object_grant", "object_type", "object_name", "priviles", "role", "grant_option", "shoprunner_terraform"),
					resource.TestCheckResourceAttr("snowflake_account_object_grant", "name", "tf-test"),
				),
			},
		},
	})
}

var testSnowflakeRoleConfig = `resource "snowflake_account_object_grant" "shoprunner_terraform" {
	object_type = "test_type"
	object_name = "test_name"
	priviliges = ["privilege1", privilege2"]
	role = "test_role"
	grant_option = true
}`
