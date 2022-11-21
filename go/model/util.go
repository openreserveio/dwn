package model

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/multiformats/go-multibase"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
	"github.com/openreserveio/dwn/go/did"
)

func CreateMessage(authorDID string, recipientDID string, dataFormat string, data []byte, methodName string) *Message {

	// Verify Message Name

	// If there is data, base64 encode it in string form
	var encodedData string = ""
	if data != nil {
		encodedData = base64.URLEncoding.EncodeToString(data)
	}

	// Start the Message
	message := Message{
		Data: encodedData,
		Processing: MessageProcessing{
			Nonce:        uuid.NewString(),
			AuthorDID:    authorDID,
			RecipientDID: recipientDID,
		},
	}

	// create the descriptor
	var dataCID string = ""
	if message.Data != "" {
		dataCID = CreateDataCID(message.Data)
	}

	messageDesc := Descriptor{
		Nonce:      uuid.New().String(),
		Method:     methodName,
		DataCID:    dataCID,
		DataFormat: dataFormat,
	}
	message.Descriptor = messageDesc

	return &message

}

func CreateDataCID(data string) string {

	d, err := qp.BuildList(basicnode.Prototype.Any, 1, func(la datamodel.ListAssembler) {
		qp.ListEntry(la, qp.String(data))
	})
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	dagcbor.Encode(d, &buf)

	cidPrefix := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}
	cid, err := cidPrefix.Sum(buf.Bytes())
	if err != nil {
		return ""
	}

	return cid.String()

}

func CreateDescriptorCID(descriptor Descriptor) string {

	d, err := qp.BuildMap(basicnode.Prototype.Any, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "method", qp.String(descriptor.Method))
		qp.MapEntry(ma, "schema", qp.String(descriptor.Schema))
		qp.MapEntry(ma, "dataCid", qp.String(descriptor.DataCID))
		qp.MapEntry(ma, "nonce", qp.String(descriptor.Nonce))
		qp.MapEntry(ma, "dataFormat", qp.String(descriptor.DataFormat))
	})
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	dagcbor.Encode(d, &buf)

	cidPrefix := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}
	cid, err := cidPrefix.Sum(buf.Bytes())
	if err != nil {
		return ""
	}

	return cid.String()

}

func CreateProcessingCID(mp MessageProcessing) string {

	d, err := qp.BuildMap(basicnode.Prototype.Any, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "nonce", qp.String(mp.Nonce))
		qp.MapEntry(ma, "author", qp.String(mp.AuthorDID))
		qp.MapEntry(ma, "recipient", qp.String(mp.RecipientDID))
	})
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	dagcbor.Encode(d, &buf)

	cidPrefix := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}
	cid, err := cidPrefix.Sum(buf.Bytes())
	if err != nil {
		return ""
	}

	return cid.String()

}

func VerifyAttestation(message *Message) bool {

	encodedProtectedHeader := message.Attestation.Signatures[0].Protected
	encodedSignature := message.Attestation.Signatures[0].Signature
	encodedPayload := message.Attestation.Payload

	// Make sure the payloads match as expected
	expectedPayload := map[string]string{
		"descriptorCid": CreateDescriptorCID(message.Descriptor),
		"processingCid": CreateProcessingCID(message.Processing),
	}
	expectedJsonPayload, _ := json.Marshal(&expectedPayload)
	expectedJwsPayload := base64.URLEncoding.EncodeToString(expectedJsonPayload)

	if expectedJwsPayload != encodedPayload {
		// These should match
		return false
	}

	// Get the ecdsa.PublicKey
	jsonProtectedHeader, err := base64.URLEncoding.DecodeString(encodedProtectedHeader)
	if err != nil {
		return false
	}

	var protectedHeaderMap map[string]string
	err = json.Unmarshal(jsonProtectedHeader, &protectedHeaderMap)
	if err != nil {
		return false
	}

	attestorDid := protectedHeaderMap["kid"]
	if attestorDid == "" {
		return false
	}

	res := did.ResolvePublicKey(attestorDid)
	if res == nil {
		return false
	}

	var publicKey *ecdsa.PublicKey = res.(*ecdsa.PublicKey)
	err = jwt.SigningMethodES512.Verify(fmt.Sprintf("%s.%s", encodedProtectedHeader, encodedPayload), encodedSignature, publicKey)
	if err != nil {
		return false
	}

	return true
}

func CreateAttestation(message *Message, privateKey ecdsa.PrivateKey) DWNJWS {

	attestation := DWNJWS{}

	// PEM encode public key -> multibase(base64url)
	publicKey := privateKey.PublicKey
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(&publicKey)
	pemEncodedPublicKey := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC EC KEY", Bytes: publicKeyBytes})
	publicKeyMultibase, _ := multibase.Encode(multibase.Base64url, pemEncodedPublicKey)

	protectedHeader := map[string]string{
		"typ": "JWS",
		"alg": jwt.SigningMethodES512.Alg(),
		"kid": fmt.Sprintf("did:key:%s", publicKeyMultibase),
	}
	jsonProtectedHeader, _ := json.Marshal(&protectedHeader)
	jwsProtectedHeader := base64.URLEncoding.EncodeToString(jsonProtectedHeader)

	payload := map[string]string{
		"descriptorCid": CreateDescriptorCID(message.Descriptor),
		"processingCid": CreateProcessingCID(message.Processing),
	}
	jsonPayload, _ := json.Marshal(&payload)
	jwsPayload := base64.URLEncoding.EncodeToString(jsonPayload)

	var jwsPayloadBytes []byte = make([]byte, base64.URLEncoding.EncodedLen(len(jsonPayload)))
	base64.URLEncoding.Encode(jwsPayloadBytes, jsonPayload)

	sig, err := jwt.SigningMethodES512.Sign(fmt.Sprintf("%s.%s", jwsProtectedHeader, jwsPayload), &privateKey)
	if err != nil {

	}

	attestation.Payload = jwsPayload
	attestation.Signatures = []DWNJWSSig{
		{
			Protected: jwsProtectedHeader,
			Signature: sig,
		},
	}

	return attestation

}
