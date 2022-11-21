package model_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/model"
)

var _ = Describe("Util", func() {

	Describe("CreateMessage", func() {

		Context("Without data", func() {

			message := model.CreateMessage("did:tmp:1", "did:tmp:2", "", nil, "CollectionsQuery", "")

			It("Should be a valid and well formed message", func() {

				Expect(message).ToNot(BeNil())

				Expect(message.Processing.Nonce).ToNot(BeEmpty())
				Expect(message.Processing.AuthorDID).To(Equal("did:tmp:1"))
				Expect(message.Processing.RecipientDID).To(Equal("did:tmp:2"))

				Expect(message.Descriptor.Nonce).ToNot(BeEmpty())
				Expect(message.Descriptor.Method).To(Equal("CollectionsQuery"))
				Expect(message.Descriptor.DataCID).To(Equal(""))
				Expect(message.Descriptor.DataFormat).To(Equal(""))

			})

		})

		Context("With data", func() {

			message := model.CreateMessage("did:tmp:10", "did:tmp:20", model.DATA_FORMAT_JSON, []byte("{\"name\":\"test\"}"), "CollectionsWrite", "")
			decodedData, err := base64.URLEncoding.DecodeString(message.Data)

			It("The Data should be decoded and match what was passed in", func() {

				Expect(err).To(BeNil())
				Expect(decodedData).ToNot(BeNil())
				Expect(decodedData).To(Equal([]byte("{\"name\":\"test\"}")))

			})

			It("Should be a valid and well formed message with proper CIDs", func() {

				Expect(message).ToNot(BeNil())

				Expect(message.Processing.Nonce).ToNot(BeEmpty())
				Expect(message.Processing.AuthorDID).To(Equal("did:tmp:10"))
				Expect(message.Processing.RecipientDID).To(Equal("did:tmp:20"))

				Expect(message.Descriptor.Nonce).ToNot(BeEmpty())
				Expect(message.Descriptor.Method).To(Equal("CollectionsWrite"))
				Expect(message.Descriptor.DataCID).ToNot(BeEmpty())
				Expect(message.Descriptor.DataFormat).To(Equal(model.DATA_FORMAT_JSON))

			})

		})

	})

	Describe("Messages with attestations", func() {

		Context("With correct signature, should verify", func() {

			It("Should verify", func() {

				data := "{\"name\":\"test user\"}"
				message := model.CreateMessage("did:tmp:1", "did:tmp:2", "", []byte(data), "CollectionsWrite", "")
				privateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

				attestation := model.CreateAttestation(message, *privateKey)
				message.Attestation = attestation

				result := model.VerifyAttestation(message)
				Expect(result).To(BeTrue())

			})

		})

		Context("With bad signature, should not verify", func() {

			It("Should not verify", func() {

				data := "{\"name\":\"test user\"}"
				message := model.CreateMessage("did:tmp:1", "did:tmp:2", "", []byte(data), "CollectionsWrite", "")
				privateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

				attestation := model.CreateAttestation(message, *privateKey)
				attestation.Signatures[0].Signature = "12345"
				message.Attestation = attestation

				result := model.VerifyAttestation(message)
				Expect(result).To(BeFalse())

			})

		})

	})

})
