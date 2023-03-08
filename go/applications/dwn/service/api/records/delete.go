package records

import (
	"context"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"net/http"
)

func RecordsDelete(ctx context.Context, collSvcClient services.RecordServiceClient, message *model.Message) model.MessageResultObject {

	// Instrumentation
	ctx, childSpan := observability.Tracer.Start(ctx, "RecordsDelete")
	defer childSpan.End()

	var messageResultObj model.MessageResultObject

	messageResultObj.Status = model.ResponseStatus{Code: http.StatusMethodNotAllowed, Detail: "CollectionDelete for next iteration"}

	return messageResultObj

}
