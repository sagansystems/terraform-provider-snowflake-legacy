package main

import (
	"github.com/sagansystems/terraform-provider-snowflake/snowflake"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: snowflake.Provider})
}
