package main

import (
	"github.com/hashicorp/terraform/plugin"
	"provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
