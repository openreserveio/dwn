package collections_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/integration-tests/testutils"
	"net/http"
	"time"
)

var _ = Describe("Write Collection", func() {

	Describe("Writing a brand new collection item", func() {

		authorPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		authorPublicKey := authorPrivateKey.PublicKey
		authorDID, _ := testutils.CreateKeyDID(&authorPublicKey)

		recipientPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		recipientPublicKey := recipientPrivateKey.PublicKey
		recipientDID, _ := testutils.CreateKeyDID(&recipientPublicKey)

		body := []byte("{\"name\":\"test\"}")
		bodyEncoded := base64.RawURLEncoding.EncodeToString(body)

		It("Stores the message correctly", func() {

			descriptor := model.Descriptor{
				Method:          model.METHOD_COLLECTIONS_WRITE,
				DataCID:         model.CreateDataCID(bodyEncoded),
				DataFormat:      model.DATA_FORMAT_JSON,
				ParentID:        "",
				Protocol:        "",
				ProtocolVersion: "",
				Schema:          "https://openreserve.io/schemas/test.json",
				CommitStrategy:  "",
				Published:       false,
				DateCreated:     time.Now(),
				DatePublished:   time.Now(),
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

			ro := model.RequestObject{}
			ro.Messages = append(ro.Messages, message)

			res, err := resty.New().R().
				SetBody(ro).
				SetHeader("Content-Type", "application/json").
				Post("http://localhost:8080/")

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

	})

})
