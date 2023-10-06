package main

import (
	"context"
	"github.com/openreserveio/dwn/go/applications/dwn/cmd"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/observability"
	"os"
)

func main() {

	ctx := context.Background()
	serviceName := "DWN"
	sd, err := observability.InitProviderWithOTELExporter(ctx, serviceName)
	if err != nil {
		log.Fatal("Unable to init tracing module:  %v", err)
		os.Exit(1)
	}
	defer sd(ctx)

	cmd.Execute(ctx)

}
