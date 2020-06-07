package main

import (
	"github.com/jerryjvl/terraform-provider-buildkite/buildkite"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return buildkite.Provider()
		},
	})
}
