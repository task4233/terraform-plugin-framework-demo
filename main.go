package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/task4233/terraform-plugin-framework-demo/log"
)

func main() {
	tfsdk.Serve(context.Background(), log.NewProvider, tfsdk.ServeOpts{
		Name: "plugin",
	})
}
