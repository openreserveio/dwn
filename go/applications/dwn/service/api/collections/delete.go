package collections

import (
	"context"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"net/http"
)

func CollectionsDelete(ctx context.Context, collSvcClient services.CollectionServiceClient, message *model.Message) model.MessageResultObject {

	// Instrumentation
	_, childSpan := observability.Tracer.Start(ctx, "CollectionsDelete")
	defer childSpan.End()

	var messageResultObj model.MessageResultObject

	messageResultObj.Status = model.ResponseStatus{Code: http.StatusMethodNotAllowed, Detail: "CollectionDelete for next iteration"}

	return messageResultObj

}
