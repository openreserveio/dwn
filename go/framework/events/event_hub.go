package events

import (
	"errors"
	"github.com/nats-io/nats.go"
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
