package collections_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/client"
	"os"
)

var _ = Describe("Write Collection", func() {

	DWN_HOST := os.Getenv("DWN_API_HOST")
	DWN_PORT := os.Getenv("DWN_API_PORT")
	TEST_SCHEMA := "https://openreserve.io/schemas/test.json"

	Describe("Writing a brand new collection item", func() {

		dwnClient := client.CreateDWNClient(fmt.Sprintf("http://%s:%s/", DWN_HOST, DWN_PORT))

		authorKeypair := client.New()
		authorIdentity := client.FromKeypair(authorKeypair)

		recipientKeypair := client.New()
		recipientIdentity := client.FromKeypair(recipientKeypair)

		body := []byte("{\"name\":\"test\"}")

		It("Stores the message correctly as its initial entry", func() {

			recordId, err := dwnClient.SaveData(TEST_SCHEMA, body, client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &authorIdentity, &recipientIdentity)

			Expect(err).To(BeNil())
			Expect(recordId).ToNot(BeEmpty())
			println(fmt.Sprintf("Record ID:  %s", recordId))

		})

		It("Creates a new entry, then updates it", func() {

			body := []byte("{\"name\":\"test_two\", \"status\":\"APPROVED\"}")
			recordId, err := dwnClient.SaveData(TEST_SCHEMA, body, client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &authorIdentity, &recipientIdentity)
			Expect(err).To(BeNil())
			Expect(recordId).ToNot(BeEmpty())

			bodyUpdated := []byte("{\"name\":\"test_two_changed\", \"status\":\"APPROVED_changed\"}")
			err = dwnClient.UpdateData(TEST_SCHEMA, recordId, bodyUpdated, client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &recipientIdentity)
			Expect(err).To(BeNil())

		})

	})

})
