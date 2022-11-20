package collsvc

import (
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/storage"
	"github.com/openreserveio/dwn/go/storage/docdbstore"
)

type CollectionService struct {
	services.UnimplementedCollectionServiceServer
	CollectionStore storage.CollectionStore
}

func CreateCollectionService(collectionStoreConnectionURI string) (*CollectionService, error) {

	// Setup Collection Store
	colLStore, err := docdbstore.CreateCollectionDocumentDBStore(collectionStoreConnectionURI)
	if err != nil {
		log.Fatal("Unable to connect to collections store:  %v", err)
		return nil, err
	}

	collService := CollectionService{
		CollectionStore: colLStore,
	}

	return &collService, nil

}
