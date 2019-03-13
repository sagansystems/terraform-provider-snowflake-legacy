package snowflake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccRoleGrantSnowflakeDatabase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testSnowflakeProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnowflakeRoleGrantConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"snowflake_role_grant", "name", "test"),
					resource.TestCheckResourceAttr("snowflake_role_grant", "role", "tf-test"),
					resource.TestCheckResourceAttr("snowflake_role_grant", "user", "tf-test-user"),
				),
			},
		},
	})
}

var testSnowflakeRoleGrantConfig = `
resource "snowflake_role_grant" "test" {
  role = "tf-test-role"
  user = "tf-test-user"
}
`
