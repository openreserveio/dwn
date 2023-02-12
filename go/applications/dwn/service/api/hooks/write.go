package hooks

import (
	"context"
	"github.com/openreserveio/dwn/go/generated/services"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"net/http"
)

func HooksWrite(ctx context.Context, hookServiceClient services.HookServiceClient, message *model.Message) model.MessageResultObject {

	// Instrumentation
	_, childSpan := observability.Tracer.Start(ctx, "CollectionsCommit")
	defer childSpan.End()

	var messageResultObj model.MessageResultObject

	// First, make sure attestations are valid and correct for this message
	// TODO:  Deal with whitelisting, blacklisting, authentication requirements
	if !model.VerifyAttestation(message) {
		messageResultObj.Status = model.ResponseStatus{Code: http.StatusUnauthorized, Detail: "Unable to verify attestation(s)."}
		return messageResultObj
	}

	// Make sure authorizations are valid for messages that are writes to existing records

	// See if there's an existing Hook record and be sure we're authorized
	// otherwise, create a new one

	return messageResultObj

}
