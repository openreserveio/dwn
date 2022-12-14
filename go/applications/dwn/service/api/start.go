package api

import (
	"github.com/openreserveio/dwn/go/applications/dwn/configuration"
	"github.com/openreserveio/dwn/go/framework"
	"github.com/openreserveio/dwn/go/log"
)

func Start(config configuration.Configuration) error {

	apiServiceOptions := framework.ServiceOptions{
		Address:    config.GetAPIListenAddress(),
		Port:       config.GetAPIListenPort(),
		SecureFlag: false,
	}

	collSvcOptions := framework.ServiceOptions{
		Address:    config.GetCollectionServiceExternalAddress(),
		Port:       config.GetCollectionServiceExternalPort(),
		SecureFlag: false,
	}

	apiService, err := CreateAPIService(&apiServiceOptions, &collSvcOptions)
	if err != nil {
		log.Fatal("Unable to create API Service:  %v", err)
		return err
	}

	return apiService.Run()

}
