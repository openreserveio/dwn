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
)

func VerifyAuthorization(message *Message) bool {

	//if message.Authorization.Payload == "" || len(message.Authorization.Signatures) == 0 {
	//	return false
	//}
	//
	//encodedProtectedHeader := message.Authorization.Signatures[0].Protected
	//encodedSignature := message.Authorization.Signatures[0].Signature
	//encodedPayload := message.Authorization.Payload
	//
	//// Get the ecdsa.PublicKey
	//jsonProtectedHeader, err := base64.URLEncoding.DecodeString(encodedProtectedHeader)
	//if err != nil {
	//	return false
	//}
	//
	//var protectedHeaderMap map[string]string
	//err = json.Unmarshal(jsonProtectedHeader, &protectedHeaderMap)
	//if err != nil {
	//	return false
	//}
	//
	//authorizerDid := protectedHeaderMap["kid"]
	//if authorizerDid == "" {
	//	return false
	//}
	//
	//res := didder.ResolvePublicKey(authorizerDid)
	//if res == nil {
	//	return false
	//}
	//
	//var publicKey *ecdsa.PublicKey = res.(*ecdsa.PublicKey)
	//err = jwt.SigningMethodES512.Verify(fmt.Sprintf("%s.%s", encodedProtectedHeader, encodedPayload), encodedSignature, publicKey)
	//if err != nil {
	//	return false
	//}

	return true
}

func CreateAuthorization(message *Message, privateKey ecdsa.PrivateKey) DWNJWS {

	authorization := DWNJWS{}

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

	sig, _ := jwt.SigningMethodES512.Sign(fmt.Sprintf("%s.%s", jwsProtectedHeader, jwsPayload), &privateKey)

	authorization.Payload = jwsPayload
	authorization.Signatures = []DWNJWSSig{
		{
			Protected: jwsProtectedHeader,
			Signature: sig,
		},
	}

	return authorization

}
