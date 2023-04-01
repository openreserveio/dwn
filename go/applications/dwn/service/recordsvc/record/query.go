package record

import (
	"context"
	"fmt"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/openreserveio/dwn/go/storage"
)

func RecordQuery(ctx context.Context, recordStore storage.RecordStore, queryMessage *model.Message) (*model.Message, error) {

	// tracing
	ctx, sp := observability.Tracer.Start(ctx, "recordsvc.record.RecordQuery")
	defer sp.End()

	var responseMessage model.Message

	// Inspect the message, and see how we should query for the record
	if queryMessage.Descriptor.Filter.RecordID != "" && queryMessage.Descriptor.Filter.Schema != "" {

		sp.AddEvent(fmt.Sprintf("Querying for record by ID %s and Schema %s", queryMessage.Descriptor.Filter.RecordID, queryMessage.Descriptor.Filter.Schema))
		latestRecord := recordStore.GetRecord(queryMessage.Descriptor.Filter.RecordID)
		if latestRecord == nil {
			sp.AddEvent("Unable to find record by ID and Schema")
			return nil, nil
		}

		sp.AddEvent(fmt.Sprintf("Querying for latest message entry by ID %s", latestRecord.LatestEntryID))
		latestMessageEntry := recordStore.GetMessageEntryByID(latestRecord.LatestEntryID)
		if latestMessageEntry == nil {
			sp.AddEvent("Unable to find latest message entry by ID")
			return nil, nil
		}

		responseMessage = latestMessageEntry.Message

	}

	return &responseMessage, nil
}
