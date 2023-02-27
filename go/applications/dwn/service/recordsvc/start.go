package recordsvc

import (
	"context"
	"fmt"
	"github.com/openreserveio/dwn/go/applications/dwn/configuration"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"google.golang.org/grpc"
	"net"
	"os"
)

func Start(ctx context.Context, config configuration.Configuration) error {

	// Start Record Service
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GetCollectionServiceListenAddress(), config.GetCollectionServiceListenPort()))
	if err != nil {
		log.Fatal("Unable to listen to address and port:  %v", err)
		os.Exit(1)
	}

	recordService, err := CreateRecordService(config.GetCollectionServiceDocumentDBURI())
	if err != nil {
		log.Fatal("Unable to start Collection Service:  %v", err)
		os.Exit(1)
	}

	// GRPC
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	services.RegisterCollectionServiceServer(grpcServer, recordService)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("GRPC Server Failed:  %v", err)
		return err
	}

	return nil

}
