package collections_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/integration-tests/testutils"
	"net/http"
)

var _ = Describe("Query For A Collection", func() {

	authorPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	authorPublicKey := authorPrivateKey.PublicKey
	authorDID, _ := testutils.CreateKeyDID(&authorPublicKey)

	recipientPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	recipientPublicKey := recipientPrivateKey.PublicKey
	recipientDID, _ := testutils.CreateKeyDID(&recipientPublicKey)

	Describe("Query for a collection that doesn't exist", func() {

		descriptor := model.Descriptor{
			Method: model.METHOD_COLLECTIONS_QUERY,
			Filter: model.CollectionsQueryFilter{
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
				Post("http://localhost:8080/")

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

})
