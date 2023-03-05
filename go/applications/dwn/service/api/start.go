package api

import (
	"context"
	"github.com/openreserveio/dwn/go/applications/dwn/configuration"
	"github.com/openreserveio/dwn/go/framework"
	"github.com/openreserveio/dwn/go/log"
)

func Start(ctx context.Context, config configuration.Configuration) error {

	apiServiceOptions := framework.ServiceOptions{
		Address:    config.GetAPIListenAddress(),
		Port:       config.GetAPIListenPort(),
		SecureFlag: false,
	}

	recordSvcOptions := framework.ServiceOptions{
		Address:    config.GetCollectionServiceExternalAddress(),
		Port:       config.GetCollectionServiceExternalPort(),
		SecureFlag: false,
	}

	hookSvcOptions := framework.ServiceOptions{
		Address:    config.GetHookServiceExternalAddress(),
		Port:       config.GetHookServiceExternalPort(),
		SecureFlag: false,
	}

	apiService, err := CreateAPIService(&apiServiceOptions, &recordSvcOptions, &hookSvcOptions)
	if err != nil {
		log.Fatal("Unable to create API Service:  %v", err)
		return err
	}

	return apiService.Run()

}
