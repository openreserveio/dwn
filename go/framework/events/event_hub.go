package events

import (
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/openreserveio/dwn/go/generated/events"
	"github.com/openreserveio/dwn/go/log"
	"google.golang.org/protobuf/proto"
)

type EventHub struct {
	Client *nats.Conn
	Subs   map[string]*nats.Subscription
}

func CreateEventHub(natsUrl string) (*EventHub, error) {

	eh := EventHub{
		Subs: make(map[string]*nats.Subscription),
	}

	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}
	eh.Client = nc

	return &eh, nil

}

func (eh *EventHub) Publish(queueName string, messageBody []byte) error {
	return eh.Client.Publish(queueName, messageBody)
}

func (eh *EventHub) Subscribe(queueName string) chan *nats.Msg {

	messagesChan := make(chan *nats.Msg)
	sub, _ := eh.Client.ChanSubscribe(queueName, messagesChan)
	eh.Subs[queueName] = sub

	return messagesChan

}

func (eh *EventHub) Unsubscribe(queueName string) error {
	if eh.Subs[queueName] == nil {
		return errors.New("No subscription")
	}
	return eh.Subs[queueName].Unsubscribe()
}

func (eh *EventHub) RaiseCreateRecordEvent(recordId string) {

	event := events.Event{
		EventType: events.EventType_CREATE_RECORD,
		RecordId:  recordId,
	}
	encodedEvent := eh.EncodeEventMessage(event)
	eh.Publish(CreateRecordEventQueue, encodedEvent)

}

func (eh *EventHub) RaiseNotifyCallbackHTTP(schemaUrl string, recordId string, protocol string, protocolVersion string, callbackUrl string) {

	discData := map[string]string{DISC_DATA_KEY_CALLBACK_URI: callbackUrl}
	event := events.Event{
		EventType:              events.EventType_NOTIFY_CALLBACK_HTTP,
		RecordId:               recordId,
		Schema:                 schemaUrl,
		Protocol:               protocol,
		ProtocolVersion:        protocolVersion,
		EventDiscretionaryData: discData,
	}
	encodedEvent := eh.EncodeEventMessage(event)
	eh.Publish(NotifyCallbackQueue, encodedEvent)

}

func (eh *EventHub) EncodeEventMessage(eventMessage events.Event) []byte {
	encodedEvent, _ := proto.Marshal(&eventMessage)
	return encodedEvent
}

func (eh *EventHub) DecodeEventMessage(encodedMessage []byte) *events.Event {

	var eventMessage events.Event
	err := proto.Unmarshal(encodedMessage, &eventMessage)
	if err != nil {
		log.Error("Error while unmarshalling event message:  %v", err)
		return nil
	}

	return &eventMessage

}
