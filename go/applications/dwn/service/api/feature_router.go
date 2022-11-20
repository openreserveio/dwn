package api

import (
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
)

type FeatureRouter struct {
	model.FeatureDetection
	CollectionServiceClient services.CollectionServiceClient
}

func CreateFeatureRouter(collsvcClient services.CollectionServiceClient) (*FeatureRouter, error) {

	return &FeatureRouter{
		FeatureDetection:        model.CurrentFeatureDetection,
		CollectionServiceClient: collsvcClient,
	}, nil

}

func (fr *FeatureRouter) Route(requestObject *model.RequestObject) (interface{}, error) {

	return nil, nil

}
