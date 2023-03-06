package records_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/go/storage"
	"github.com/openreserveio/dwn/integration-tests/testutils"
	"net/http"
	"os"
	"time"
)

var _ = Describe("Query For A Record", func() {

	DWN_HOST := os.Getenv("DWN_API_HOST")
	DWN_PORT := os.Getenv("DWN_API_PORT")

	authorPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	authorPublicKey := authorPrivateKey.PublicKey
	authorDID, _ := testutils.CreateKeyDID(&authorPublicKey)

	recipientPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	recipientPublicKey := recipientPrivateKey.PublicKey
	recipientDID, _ := testutils.CreateKeyDID(&recipientPublicKey)

	Describe("Query for a collection that doesn't exist", func() {

		descriptor := model.Descriptor{
			Method: model.METHOD_RECORDS_QUERY,
			Filter: model.DescriptorFilter{
				RecordID: "DOES NOT EXIST",
				Schema:   "https://openreserve.io/schemas/test.json",
			},
		}

		messageProcessing := model.MessageProcessing{
			Nonce:        uuid.NewString(),
			AuthorDID:    authorDID,
			RecipientDID: recipientDID,
		}

		message := model.Message{
			ContextID:  "",
			Processing: messageProcessing,
			Descriptor: descriptor,
		}

		authorization := model.CreateAuthorization(&message, *authorPrivateKey)
		attestation := model.CreateAttestation(&message, *authorPrivateKey)
		message.Attestation = attestation
		message.Authorization = authorization

		It("Queries the DWN for an entry that does not exist", func() {

			ro := model.RequestObject{}
			ro.Messages = append(ro.Messages, message)

			res, err := resty.New().R().
				SetBody(ro).
				SetHeader("Content-Type", "application/json").
				Post(fmt.Sprintf("http://%s:%s/", DWN_HOST, DWN_PORT))

			Expect(err).To(BeNil())
			Expect(res).ToNot(BeNil())
			Expect(res.StatusCode()).To(Equal(http.StatusOK))

			var responseObject model.ResponseObject
			err = json.Unmarshal(res.Body(), &responseObject)
			Expect(err).To(BeNil())

			Expect(responseObject.Status.Code).To(Equal(http.StatusOK))
			Expect(len(responseObject.Replies)).To(Equal(1))
			Expect(responseObject.Replies[0].Status.Code).To(Equal(http.StatusNotFound))

		})

	})

	Describe("Query for a Record that does exist", func() {

		// Need to Create the Collection Record.
		// TODO:  Refactor using new client lib to abstract this
		body := []byte("{\"name\":\"test\"}")
		bodyEncoded := base64.RawURLEncoding.EncodeToString(body)

		descriptor := model.Descriptor{
			Method:          model.METHOD_RECORDS_WRITE,
			DataCID:         model.CreateDataCID(bodyEncoded),
			DataFormat:      model.DATA_FORMAT_JSON,
			ParentID:        "",
			Protocol:        "",
			ProtocolVersion: "",
			Schema:          "https://openreserve.io/schemas/test.json",
			CommitStrategy:  "",
			Published:       false,
			DateCreated:     time.Now(),
			DatePublished:   nil,
		}

		processing := model.MessageProcessing{
			Nonce:        uuid.NewString(),
			AuthorDID:    authorDID,
			RecipientDID: recipientDID,
		}

		descriptorCID := model.CreateDescriptorCID(descriptor)
		processingCID := model.CreateProcessingCID(processing)
		recordId := model.CreateRecordCID(descriptorCID, processingCID)

		message := model.Message{
			RecordID:   recordId,
			ContextID:  "",
			Data:       bodyEncoded,
			Processing: processing,
			Descriptor: descriptor,
		}

		attestation := model.CreateAttestation(&message, *authorPrivateKey)
		message.Attestation = attestation

		authorization := model.CreateAuthorization(&message, *authorPrivateKey)
		message.Authorization = authorization

		It("Stores the message as the initial entry", func() {

			ro := model.RequestObject{}
			ro.Messages = append(ro.Messages, message)

			res, err := resty.New().R().
				SetBody(ro).
				SetHeader("Content-Type", "application/json").
				Post(fmt.Sprintf("http://%s:%s/", DWN_HOST, DWN_PORT))

			Expect(err).To(BeNil())
			Expect(res).ToNot(BeNil())
			Expect(res.StatusCode()).To(Equal(http.StatusOK))

			var responseObject model.ResponseObject
			err = json.Unmarshal(res.Body(), &responseObject)
			Expect(err).To(BeNil())

			Expect(responseObject.Status.Code).To(Equal(http.StatusOK))
			Expect(len(responseObject.Replies)).To(Equal(1))
			Expect(responseObject.Replies[0].Status.Code).To(Equal(http.StatusOK))
			Expect(responseObject.Replies[0].Entries[0].Result).ToNot(BeNil())

		})

		It("Queries for the record just created", func() {

			// TODO:  THis should be refactored into client lib
			queryDescriptor := model.Descriptor{
				Method: model.METHOD_RECORDS_QUERY,
				Filter: model.DescriptorFilter{
					RecordID: message.RecordID,
					Schema:   message.Descriptor.Schema,
				},
			}

			queryMessageProcessing := model.MessageProcessing{
				Nonce:        uuid.NewString(),
				AuthorDID:    authorDID,
				RecipientDID: recipientDID,
			}

			queryMessage := model.Message{
				ContextID:  "",
				Processing: queryMessageProcessing,
				Descriptor: queryDescriptor,
			}

			authorization := model.CreateAuthorization(&queryMessage, *authorPrivateKey)
			attestation := model.CreateAttestation(&queryMessage, *authorPrivateKey)
			queryMessage.Attestation = attestation
			queryMessage.Authorization = authorization

			ro := model.RequestObject{}
			ro.Messages = append(ro.Messages, queryMessage)

			res, err := resty.New().R().
				SetBody(ro).
				SetHeader("Content-Type", "application/json").
				Post(fmt.Sprintf("http://%s:%s/", DWN_HOST, DWN_PORT))

			Expect(err).To(BeNil())
			Expect(res).ToNot(BeNil())
			Expect(res.StatusCode()).To(Equal(http.StatusOK))

			var responseObject model.ResponseObject
			err = json.Unmarshal(res.Body(), &responseObject)
			Expect(err).To(BeNil())

			Expect(responseObject.Status.Code).To(Equal(http.StatusOK))
			Expect(len(responseObject.Replies)).To(Equal(1))
			Expect(responseObject.Replies[0].Status.Code).To(Equal(http.StatusOK))

			// JSON Decode
			// TODO: Change this return object -- shouldn't be a message entry from storage package
			var entry storage.MessageEntry
			err = json.Unmarshal(responseObject.Replies[0].Entries[0].Result, &entry)
			Expect(err).To(BeNil())
			Expect(entry.RecordID).To(Equal(message.RecordID))

		})

	})

})
