package model_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/client"
	"github.com/openreserveio/dwn/go/model"
)

var _ = Describe("Util", func() {

	Describe("CreateMessage", func() {

		Context("Without data", func() {

			message := model.CreateMessage("did:tmp:1", "did:tmp:2", "", nil, "Records", "Query", "", "https://openreserve.io/schemas/test.json")

			It("Should be a valid and well formed message", func() {

				Expect(message).ToNot(BeNil())

				Expect(message.Processing.Nonce).ToNot(BeEmpty())
				Expect(message.Processing.AuthorDID).To(Equal("did:tmp:1"))
				Expect(message.Processing.RecipientDID).To(Equal("did:tmp:2"))

				Expect(message.Descriptor.Method).To(Equal("Query"))
				Expect(message.Descriptor.DataCID).To(Equal(""))
				Expect(message.Descriptor.DataFormat).To(Equal(""))

			})

		})

		Context("With data", func() {

			protocolDefinition := model.ProtocolDefinition{
				ContextID:       "",
				Protocol:        "",
				ProtocolVersion: "",
			}

			authorKeypair := client.NewKeypair()
			authorID := client.FromKeypair(authorKeypair)
			recipKeypair := client.NewKeypair()
			recipID := client.FromKeypair(recipKeypair)

			message := model.CreateInitialRecordsWriteMessage(authorID.DID, recipID.DID, &protocolDefinition, "https://openreserve.io/schemas/test.json", model.DATA_FORMAT_JSON, []byte("{\"name\":\"test\"}"))
			decodedData, err := base64.URLEncoding.DecodeString(message.Data)
			println(fmt.Sprintf("Message Record ID:  %v", message.RecordID))

			It("The Data should be decoded and match what was passed in", func() {

				Expect(err).To(BeNil())
				Expect(decodedData).ToNot(BeNil())
				Expect(decodedData).To(Equal([]byte("{\"name\":\"test\"}")))

			})

			It("Should be a valid and well formed message with proper CIDs", func() {

				Expect(message).ToNot(BeNil())

				Expect(message.Processing.Nonce).ToNot(BeEmpty())
				Expect(message.Processing.AuthorDID).To(Equal(authorID.DID))
				Expect(message.Processing.RecipientDID).To(Equal(recipID.DID))

				Expect(message.Descriptor.Method).To(Equal(model.METHOD_RECORDS_WRITE))
				Expect(message.Descriptor.DataCID).ToNot(BeEmpty())
				Expect(message.Descriptor.DataFormat).To(Equal(model.DATA_FORMAT_JSON))

			})

		})

	})

	Describe("Messages with attestations", func() {

		Context("With correct signature, should verify", func() {

			It("Should verify", func() {

				publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
				authorDID, _ := did.CreateDIDKey(crypto.Ed25519, publicKey)
				data := "{\"name\":\"test user\"}"
				message := model.CreateMessage(authorDID.String(), "did:tmp:2", "", []byte(data), "CollectionsWrite", "", "", "https://openreserve.io/schemas/test.json")

				attestation := model.CreateAttestation(message, nil, privateKey)
				message.Attestation = attestation

				result := model.VerifyAttestation(message)
				Expect(result).To(BeTrue())

			})

		})

		Context("With bad signature, should not verify", func() {

			It("Should not verify", func() {

				data := "{\"name\":\"test user\"}"
				message := model.CreateMessage("did:tmp:1", "did:tmp:2", "", []byte(data), "CollectionsWrite", "", "", "https://openreserve.io/schemas/test.json")
				privateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

				attestation := model.CreateAttestation(message, nil, *privateKey)
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
