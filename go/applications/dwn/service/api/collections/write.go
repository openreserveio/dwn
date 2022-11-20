package collections

import (
	"github.com/openreserveio/dwn/go/model"
	"net/http"
)

func CollectionsWrite(message *model.Message) model.MessageResultObject {

	var messageResultObj model.MessageResultObject

	// First, find the schema and make sure it has been registered
	schemaUri := message.Descriptor.Schema
	if schemaUri == "" {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusBadRequest, Detail: "Schema URI is required for a CollectionsWrite"}
		return messageResultObj
	}

	return messageResultObj

}
