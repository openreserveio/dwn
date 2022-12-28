package model

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/multiformats/go-multibase"
	"github.com/openreserveio/dwn/go/did"
)

func VerifyAttestation(message *Message) bool {

	if message.Attestation.Payload == "" || len(message.Attestation.Signatures) == 0 {
		return false
	}

	encodedProtectedHeader := message.Attestation.Signatures[0].Protected
	encodedSignature := message.Attestation.Signatures[0].Signature
	encodedPayload := message.Attestation.Payload

	// Make sure the payloads match as expected
	//expectedPayload := map[string]string{
	//	"descriptorCid": CreateDescriptorCID(message.Descriptor),
	//	"processingCid": CreateProcessingCID(message.Processing),
	//}
	//expectedJsonPayload, _ := json.Marshal(&expectedPayload)
	//expectedJwsPayload := base64.URLEncoding.EncodeToString(expectedJsonPayload)
	//
	//if expectedJwsPayload != encodedPayload {
	//	// These should match
	//	return false
	//}

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
