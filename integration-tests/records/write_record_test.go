package records_test

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/client"
	"os"
)

var _ = Describe("Write Record", func() {

	DWN_HOST := os.Getenv("DWN_API_HOST")
	DWN_PORT := os.Getenv("DWN_API_PORT")
	TEST_SCHEMA := "https://openreserve.io/schemas/test.json"

	Describe("Writing a brand new record", func() {

		dwnClient := client.CreateDWNClient(fmt.Sprintf("http://%s:%s/", DWN_HOST, DWN_PORT))

		authorKeypair := client.NewKeypair()
		authorIdentity := client.FromKeypair(authorKeypair)

		recipientKeypair := client.NewKeypair()
		recipientIdentity := client.FromKeypair(recipientKeypair)

		randomUUID := uuid.NewString()
		body := []byte(fmt.Sprintf("{\"name\":\"test_%s\", \"status\":\"APPROVED\"}", randomUUID))

		It("Stores the message correctly as its initial entry", func() {

			recordId, err := dwnClient.SaveData(TEST_SCHEMA, body, client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &authorIdentity, &recipientIdentity)

			Expect(err).To(BeNil())
			Expect(recordId).ToNot(BeEmpty())
			println(fmt.Sprintf("Record ID:  %s", recordId))

		})

		It("Creates a new entry, then updates it", func() {

			body := []byte(fmt.Sprintf("{\"name\":\"test_%s\", \"status\":\"APPROVED\"}", randomUUID))
			recordId, err := dwnClient.SaveData(TEST_SCHEMA, body, client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &authorIdentity, &recipientIdentity)
			Expect(err).To(BeNil())
			Expect(recordId).ToNot(BeEmpty())

			bodyUpdated := []byte(fmt.Sprintf("{\"name\":\"test_%s\", \"status\":\"APPROVED_CHANGED\"}", randomUUID))
			err = dwnClient.UpdateData(TEST_SCHEMA, recordId, bodyUpdated, client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &recipientIdentity)
			Expect(err).To(BeNil())

			//_, data, _, err := dwnClient.GetData(TEST_SCHEMA, recordId, &recipientIdentity)
			//Expect(err).To(BeNil())
			//Expect(data).ToNot(BeNil())
			//Expect(data).To(Equal(bodyUpdated))

		})

	})

})
