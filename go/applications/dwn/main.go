package main

import (
	"context"
	"fmt"
	"github.com/openreserveio/dwn/go/applications/dwn/cmd"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/tracing"
	"os"
)

func main() {

	serviceName := "DWN"
	if len(os.Args) < 2 {
		serviceName = "DWN-BASE"
	} else {
		serviceName = fmt.Sprintf("%s-%s", os.Args[0], os.Args[1])
	}

	sd, err := tracing.InitProviderWithJaegerExporter(context.Background(), serviceName)
	if err != nil {
		log.Fatal("Unable to init tracing module:  %v", err)
		os.Exit(1)
	}
	defer sd(context.Background())

	cmd.Execute()

}
