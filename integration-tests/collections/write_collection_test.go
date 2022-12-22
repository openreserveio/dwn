package collections_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/model"
	"github.com/openreserveio/dwn/integration-tests/testutils"
	"net/http"
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

		It("Stores the message correctly", func() {

			message := testutils.CreateMessage(authorDID, recipientDID, "application/json", body, "CollectionsWrite")
			message.Descriptor.Schema = "https://openreserve.io/schemas/test.json"
			attestation := testutils.CreateAttestation(message, *authorPrivateKey)
			message.Attestation = attestation

			ro := model.RequestObject{}
			ro.Messages = append(ro.Messages, *message)

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
