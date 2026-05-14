package main

import (
	"context"
	"log"

	"github.com/arcusis/terraform-provider-coolify/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev"

func main() {
	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/arcusis/coolify",
	})
	if err != nil {
		log.Fatal(err)
	}
}
