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

			secondBody := []byte("{\"name\":\"test_two\", \"status\":\"APPROVED\"}")
			secondRecordId, err := dwnClient.SaveData(TEST_SCHEMA, secondBody, client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &authorIdentity, &recipientIdentity)
			Expect(err).To(BeNil())
			Expect(secondRecordId).ToNot(BeEmpty())

		})
		//
		//	It("Allows user to commit the latest message entry", func() {
		//
		//		commitDescriptor := model.Descriptor{
		//			Method:          model.METHOD_COLLECTIONS_COMMIT,
		//			DataCID:         secondDescriptor.DataCID,
		//			DataFormat:      model.DATA_FORMAT_JSON,
		//			ParentID:        recordId,
		//			Protocol:        "",
		//			ProtocolVersion: "",
		//			Schema:          "https://openreserve.io/schemas/test.json",
		//			CommitStrategy:  "",
		//			Published:       false,
		//			DateCreated:     time.Now(),
		//			DatePublished:   nil,
		//		}
		//
		//		commitProcessing := model.MessageProcessing{
		//			Nonce:        uuid.NewString(),
		//			AuthorDID:    authorDID,
		//			RecipientDID: recipientDID,
		//		}
		//
		//		commitDescriptorCID := model.CreateDescriptorCID(commitDescriptor)
		//		commitProcessingCID := model.CreateProcessingCID(commitProcessing)
		//		commitRecordId := model.CreateRecordCID(commitDescriptorCID, commitProcessingCID)
		//
		//		commitMessage := model.Message{
		//			RecordID:   commitRecordId,
		//			ContextID:  "",
		//			Processing: commitProcessing,
		//			Descriptor: commitDescriptor,
		//		}
		//
		//		commitAttestation := model.CreateAttestation(&commitMessage, *authorPrivateKey)
		//		commitMessage.Attestation = commitAttestation
		//
		//		commitAuthorization := model.CreateAuthorization(&commitMessage, *authorPrivateKey)
		//		commitMessage.Authorization = commitAuthorization
		//
		//		ro := model.RequestObject{}
		//		ro.Messages = append(ro.Messages, commitMessage)
		//
		//		res, err := resty.New().R().
		//			SetBody(ro).
		//			SetHeader("Content-Type", "application/json").
		//			Post(fmt.Sprintf("http://%s:%s/", DWN_HOST, DWN_PORT))
		//
		//		Expect(err).To(BeNil())
		//		Expect(res).ToNot(BeNil())
		//		Expect(res.StatusCode()).To(Equal(http.StatusOK))
		//
		//		var responseObject model.ResponseObject
		//		err = json.Unmarshal(res.Body(), &responseObject)
		//		Expect(err).To(BeNil())
		//
		//		Expect(responseObject.Status.Code).To(Equal(http.StatusOK))
		//		Expect(len(responseObject.Replies)).To(Equal(1))
		//		Expect(responseObject.Replies[0].Status.Code).To(Equal(http.StatusOK), fmt.Sprintf("Status: %d: %s", responseObject.Replies[0].Status.Code, responseObject.Replies[0].Status.Detail))
		//		Expect(responseObject.Replies[0].Entries[0].Result).ToNot(BeNil())
		//
		//	})

	})

})
