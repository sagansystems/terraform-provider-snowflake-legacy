package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sagansystems/terraform-provider-snowflake/snowflake"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: snowflake.Provider})
}
