package events_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/framework/events"
)

var _ = Describe("EventHub", func() {

	It("Connects a sender and receiver to nats and sends/receives a message", func() {

		eh, err := events.CreateEventHub("nats://localhost:4222")
		Expect(err).To(BeNil())
		Expect(eh).ToNot(BeNil())

		msgChan := eh.Subscribe("testing")

		err = eh.Publish("testing", []byte("hello world"))
		Expect(err).To(BeNil())

		msg := <-msgChan
		Expect(msg).ToNot(BeNil())
		Expect(msg.Data).To(Equal([]byte("hello world")))

	})

})
