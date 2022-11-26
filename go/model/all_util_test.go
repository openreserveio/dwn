package model_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/did"
	"github.com/openreserveio/dwn/go/model"
)

var _ = Describe("Util", func() {

	Describe("CreateMessage", func() {

		Context("Without data", func() {

			message := model.CreateMessage("did:tmp:1", "did:tmp:2", "", nil, "CollectionsQuery", "", "https://openreserve.io/schemas/test.json")

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

			message := model.CreateMessage("did:tmp:10", "did:tmp:20", model.DATA_FORMAT_JSON, []byte("{\"name\":\"test\"}"), "CollectionsWrite", "", "https://openreserve.io/schemas/test.json")
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

				privateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
				authorDID, _ := did.CreateKeyDID(&privateKey.PublicKey)
				data := "{\"name\":\"test user\"}"
				message := model.CreateMessage(authorDID, "did:tmp:2", "", []byte(data), "CollectionsWrite", "", "https://openreserve.io/schemas/test.json")

				attestation := model.CreateAttestation(message, *privateKey)
				message.Attestation = attestation

				result := model.VerifyAttestation(message)
				Expect(result).To(BeTrue())

			})

		})

		Context("With bad signature, should not verify", func() {

			It("Should not verify", func() {

				data := "{\"name\":\"test user\"}"
				message := model.CreateMessage("did:tmp:1", "did:tmp:2", "", []byte(data), "CollectionsWrite", "", "https://openreserve.io/schemas/test.json")
				privateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

				attestation := model.CreateAttestation(message, *privateKey)
				attestation.Signatures[0].Signature = "12345"
				message.Attestation = attestation

				result := model.VerifyAttestation(message)
				Expect(result).To(BeFalse())

			})

		})

		Context("CBOR Encoding", func() {

			node, err := qp.BuildMap(basicnode.Prototype.Any, 4, func(ma datamodel.MapAssembler) {
				qp.MapEntry(ma, "authorDID", qp.String("did:web:openreserve.io:echo"))
				qp.MapEntry(ma, "recipientDID", qp.String("did:web:did.openreserve.io:KJDSOIDH7HDFJDS8KDNCKDC8DKNKNCSD9822298KJKJDCB"))
				qp.MapEntry(ma, "content", qp.Map(2, func(maDeep datamodel.MapAssembler) {
					qp.MapEntry(maDeep, "foo", qp.String("bar"))
					qp.MapEntry(maDeep, "goo", qp.String("car"))
				}))
				qp.MapEntry(ma, "create-date", qp.String("1/1/2020"))
			})

			It("Should be a valid node", func() {
				Expect(err).To(BeNil())
				Expect(node.IsNull()).To(BeFalse())
			})

			var initialCid cid.Cid
			var secondCid cid.Cid
			It("Should encode to CBOR and CID", func() {

				var cborBuffer bytes.Buffer
				err = dagcbor.Encode(node, &cborBuffer)
				Expect(err).To(BeNil())

				cidPrefix := cid.Prefix{Version: 1}
				initialCid, err = cidPrefix.Sum(cborBuffer.Bytes())
				Expect(err).To(BeNil())
				Expect(initialCid.String()).ToNot(BeEmpty())
			})

			It("Should match another CBOR DAG ordered differently", func() {

				// Same content, different ordering
				secondNode, err := qp.BuildMap(basicnode.Prototype.Any, 4, func(ma datamodel.MapAssembler) {
					qp.MapEntry(ma, "content", qp.Map(2, func(maDeep datamodel.MapAssembler) {
						qp.MapEntry(maDeep, "foo", qp.String("bar"))
						qp.MapEntry(maDeep, "goo", qp.String("car"))
					}))
					qp.MapEntry(ma, "create-date", qp.String("1/1/2020"))
					qp.MapEntry(ma, "authorDID", qp.String("did:web:openreserve.io:echo"))
					qp.MapEntry(ma, "recipientDID", qp.String("did:web:did.openreserve.io:KJDSOIDH7HDFJDS8KDNCKDC8DKNKNCSD9822298KJKJDCB"))
				})

				var cborBuffer bytes.Buffer
				err = dagcbor.Encode(secondNode, &cborBuffer)
				Expect(err).To(BeNil())

				cidPrefix := cid.Prefix{Version: 1}
				secondCid, err = cidPrefix.Sum(cborBuffer.Bytes())
				Expect(err).To(BeNil())
				Expect(secondCid.String()).ToNot(BeEmpty())
				Expect(initialCid.Equals(secondCid)).To(BeTrue())
			})

		})

	})

})
