package procloop

import (
	"errors"
	"fmt"
	"github.com/openreserveio/dwn/go/generated/events"
)

func ProcessMessage(event *events.Event) error {

	switch event.EventType {

	case events.EventType_NOTIFY_CALLBACK_HTTP:
		return NotifyHTTP(event)

	default:
		return errors.New(fmt.Sprintf("Unknown Event Type:  %d", event.EventType))
	}

	return nil
}
