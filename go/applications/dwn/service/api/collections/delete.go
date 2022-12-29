package collections

import (
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
	"net/http"
)

func CollectionsDelete(collSvcClient services.CollectionServiceClient, message *model.Message) model.MessageResultObject {

	var messageResultObj model.MessageResultObject

	messageResultObj.Status = model.ResponseStatus{Code: http.StatusMethodNotAllowed, Detail: "CollectionDelete for next iteration"}

	return messageResultObj

}
