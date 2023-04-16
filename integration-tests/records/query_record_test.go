package records_test

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/client"
	"github.com/openreserveio/dwn/go/model"
	"os"
)

var _ = Describe("Query For A Record", func() {

	TEST_SCHEMA := "https://openreserve.io/schema/test.json"
	DWN_HOST := os.Getenv("DWN_API_HOST")
	DWN_PORT := os.Getenv("DWN_API_PORT")

	authorKeypair := client.NewKeypair()

	recipientKeypair := client.NewKeypair()
	recipientIdentity := client.FromKeypair(recipientKeypair)
	recipientDID := recipientIdentity.DID

	Describe("Query for a collection that doesn't exist", func() {

		dwnClient := client.CreateDWNClient(fmt.Sprintf("http://%s:%s/", DWN_HOST, DWN_PORT))
		message := model.CreateQueryRecordsMessage(TEST_SCHEMA, "DOES NOT EXIST", &model.ProtocolDefinition{}, recipientDID)

		authorization := model.CreateAuthorization(message, authorKeypair.PublicKey, authorKeypair.PrivateKey)
		attestation := model.CreateAttestation(message, authorKeypair.PublicKey, authorKeypair.PrivateKey)
		message.Attestation = attestation
		message.Authorization = authorization

		It("Queries the DWN for an entry that does not exist", func() {

			ro := model.RequestObject{}
			ro.Messages = append(ro.Messages, *message)

			_, data, dataFormat, err := dwnClient.GetData(TEST_SCHEMA, "DOES NOT EXIST", &recipientIdentity)

			Expect(err).To(BeNil())
			Expect(data).To(BeNil())
			Expect(dataFormat).To(BeEmpty())

		})

	})

	Describe("Query for a Record that does exist", func() {

		authorKeypair := client.NewKeypair()
		authorIdentity := client.FromKeypair(authorKeypair)

		recipientKeypair := client.NewKeypair()
		recipientIdentity := client.FromKeypair(recipientKeypair)

		var recordId string
		var err error
		dwnClient := client.CreateDWNClient(fmt.Sprintf("http://%s:%s", DWN_HOST, DWN_PORT))

		It("Stores the message as the initial entry", func() {

			// Need to Create the Record.
			recordData := []byte(fmt.Sprintf("{\"name\":\"test_%s\", \"status\":\"APPROVED\"}", uuid.NewString()))
			recordId, err = dwnClient.SaveData(TEST_SCHEMA, recordData, "application/json", &authorIdentity, &recipientIdentity)

			Expect(err).To(BeNil())
			Expect(recordId).ToNot(BeEmpty())

		})

		It("Queries for the record just created", func() {

			_, recordData, dataFormat, err := dwnClient.GetData(TEST_SCHEMA, recordId, &recipientIdentity)

			Expect(err).To(BeNil())
			Expect(recordData).ToNot(BeNil())
			Expect(dataFormat).To(Equal("application/json"))

		})

	})

})
