package hooksvc

import (
	"fmt"
	"github.com/openreserveio/dwn/go/applications/dwn/configuration"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"google.golang.org/grpc"
	"net"
	"os"
)

func Start(config configuration.Configuration) error {

	// Start Hook Service
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GetHookServiceListenAddress(), config.GetHookServiceListenPort()))
	if err != nil {
		log.Fatal("Unable to listen to address and port:  %v", err)
		os.Exit(1)
	}

	hookService, err := CreateHookService(config.GetHookServiceDocumentDBURI())
	if err != nil {
		log.Fatal("Unable to start Hook Service:  %v", err)
		os.Exit(1)
	}

	// GRPC
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	services.RegisterHookServiceServer(grpcServer, hookService)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("GRPC Server Failed:  %v", err)
		return err
	}

	return nil
}
