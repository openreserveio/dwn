package api

import (
	"fmt"
	"github.com/openreserveio/dwn/go/applications/dwn/service/api/collections"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/model"
	"net/http"
	"time"
)

type FeatureRouter struct {
	model.FeatureDetection
	MaximumTimeoutSeconds   int
	CollectionServiceClient services.CollectionServiceClient
}

func CreateFeatureRouter(collsvcClient services.CollectionServiceClient, maxTimeoutSeconds int) (*FeatureRouter, error) {

	return &FeatureRouter{
		FeatureDetection:        model.CurrentFeatureDetection,
		MaximumTimeoutSeconds:   maxTimeoutSeconds,
		CollectionServiceClient: collsvcClient,
	}, nil

}

func (fr *FeatureRouter) Route(requestObject *model.RequestObject) (interface{}, error) {

	// Setup Response Object
	responseObject := model.ResponseObject{}
	responseObject.Replies = make([]model.MessageResultObject, len(requestObject.Messages))

	// Setup chan to receive responses
	procChan := make(chan *MessageProcResult)

	// Process Messages and append responses to responseObject
	// We can probably parallel these message processors
	for idx, message := range requestObject.Messages {
		log.Info("Processing message %d", idx)
		go fr.processMessage(idx, procChan, &message)
	}

	// Listen for all responses to come back, wrap the result objects into the response object and respond
	for i := 0; i < len(requestObject.Messages); i++ {

		// Use a maximum timeout for all the message processors
		select {
		case res := <-procChan:
			log.Info("Received message response %d", res.Index)
			responseObject.Replies[res.Index] = *res.MessageResult
		case <-time.After(15 * time.Second):
			// Generic Timeout error for remaining response objects
			fr.genericTimeouts(&responseObject)
		}

	}

	responseObject.Status = model.ResponseStatus{Code: http.StatusOK}
	return &responseObject, nil

}

// Meant to be run in a goroutine
type MessageProcResult struct {
	Index         int
	MessageResult *model.MessageResultObject
}

// idx is for ordering, MessageProcResult wraps the messageresult object and the index for the responseobject
func (fr *FeatureRouter) processMessage(idx int, procComm chan *MessageProcResult, message *model.Message) {

	// Support Simple Test Messages
	if message.RecordID == "TEST" && message.Data == "TEST" {

		// This is a test message
		procComm <- &MessageProcResult{
			Index: idx,
			MessageResult: &model.MessageResultObject{
				Status:  model.ResponseStatus{Code: http.StatusOK},
				Entries: nil,
			},
		}
		return

	}

	var messageResult model.MessageResultObject

	// Route logic based on Method
	switch message.Descriptor.Method {

	case model.METHOD_COLLECTIONS_QUERY:
		messageResult = collections.CollectionsQuery(fr.CollectionServiceClient, message)

	case model.METHOD_COLLECTIONS_WRITE:
		messageResult = collections.CollectionsWrite(fr.CollectionServiceClient, message)

	case model.METHOD_COLLECTIONS_COMMIT:
		messageResult = collections.CollectionsCommit(fr.CollectionServiceClient, message)

	case model.METHOD_COLLECTIONS_DELETE:
		messageResult = collections.CollectionsDelete(fr.CollectionServiceClient, message)

	default:
		messageResult = model.MessageResultObject{Status: model.ResponseStatus{Code: http.StatusBadRequest, Detail: fmt.Sprintf("We do not yet support message method: %s", message.Descriptor.Method)}}

	}

	// Reply back to channel
	procResult := MessageProcResult{
		Index:         idx,
		MessageResult: &messageResult,
	}
	procComm <- &procResult

}

// Fills missing responses with generic timeouts
func (fr *FeatureRouter) genericTimeouts(responseObject *model.ResponseObject) {
	for idx, reply := range responseObject.Replies {
		if reply.Status.Code == 0 {
			// never got a status code of 200 or other error > 0
			responseObject.Replies[idx] = model.MessageResultObject{
				Status:  model.ResponseStatus{Code: http.StatusRequestTimeout, Detail: "This message could not be processed within a generous timeout"},
				Entries: nil,
			}
		}
	}
}
