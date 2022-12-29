package collections

import (
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
	"net/http"
)

func CollectionsCommit(collSvcClient services.CollectionServiceClient, message *model.Message) model.MessageResultObject {

	return model.MessageResultObject{
		Status:  model.ResponseStatus{Code: http.StatusOK},
		Entries: nil,
	}

}
