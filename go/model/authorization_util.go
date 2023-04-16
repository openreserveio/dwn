package model

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	didsdk "github.com/TBD54566975/ssi-sdk/did"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"strings"
)

/*
 * This entire utility file is an implementation of the DWN Message Authorization spec:
 * https://identity.foundation/decentralized-web-node/spec/#message-authorization
 */

func VerifyAuthorization(message *Message) bool {

	// Do some basic checking
	if message == nil || message.Authorization.Payload == "" || message.Authorization.Signatures == nil || len(message.Authorization.Signatures) == 0 {
		return false
	}

	jwsToVerify := fmt.Sprintf("%s.%s.%s", message.Authorization.Signatures[0].Protected, message.Authorization.Payload, message.Authorization.Signatures[0].Signature)
	jwsMessage, err := jws.ParseString(jwsToVerify)
	if err != nil {
		return false
	}

	for _, sigs := range jwsMessage.Signatures() {

		signingKeyId := sigs.ProtectedHeaders().KeyID()
		signingAlg := sigs.ProtectedHeaders().Algorithm()

		// Resolve the signing key
		signingKeyDidDocument, err := ResolveDID(signingKeyId)
		if err != nil || signingKeyDidDocument == nil {
			return false
		}

		// Get the public key from the DID Document verification method
		// put all VerificationMethods into a map
		// Get everything after the # in the signingKeyId
		ref := signingKeyId[strings.Index(signingKeyId, "#"):]
		authPublicKey, err := didsdk.GetKeyFromVerificationMethod(*signingKeyDidDocument, ref)
		if err != nil {
			return false
		}

		_, err = jws.Verify([]byte(jwsToVerify), signingAlg, authPublicKey)
		if err != nil {
			return false
		}

	}

	return true
}

func CreateAuthorization(message *Message, authKeyUri string, publicKey crypto.PublicKey, privateKey crypto.PrivateKey) DWNJWS {

	authorization := DWNJWS{}

	// PEM encode public key -> multibase(base64url)
	// publicKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	// pemEncodedPublicKey := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC EC KEY", Bytes: publicKeyBytes})
	// publicKeyMultibase, _ := multibase.Encode(multibase.Base64url, pemEncodedPublicKey)

	// Using github.com/lestrrat-go/jwx package here because *amazing*
	/*
		From Spec:
		* The JWS MUST include a payload property, and its value MUST be an object composed of the following values:
		  * The object MUST include a descriptorCid property, and its value MUST be the stringified Version 1 CID of the
			DAG CBOR encoded descriptor object.
		  * The object MAY include a permissionsGrantCid property, and its value MUST be the stringified Version 1 CID of the
			DAG CBOR encoded Permission Grant being invoked.
		  * If attestation of an object is permitted, the payload MAY include an attestationCid property, and its value MUST be
			the stringified Version 1 CID of the DAG CBOR encoded attestation string.

	*/
	descriptorCid := CreateDescriptorCID(message.Descriptor)
	// TODO:  Include permissionsGrants and attestation logic per spec
	// attestationCid := CreateAttestationCID(message.Attestation)
	// permissionsGrantCid := CreatePermissionsGrantCID(message.PermissionsGrant)

	payloadMap := map[string]string{
		"descriptorCid":       descriptorCid,
		"permissionsGrantCid": "",
		"attestationCid":      "",
	}
	jsonPayload, _ := json.Marshal(&payloadMap)

	additionalHeaders := jws.NewHeaders()
	additionalHeaders.Set("kid", authKeyUri)

	signature, err := jws.Sign(jsonPayload, jwa.EdDSA, privateKey, jws.WithHeaders(additionalHeaders))
	if err != nil {
		panic(err)
	}

	jwsMsg, err := jws.Parse(signature)
	if err != nil {
		panic(err)
	}

	authorization.Payload = base64.RawURLEncoding.EncodeToString(jwsMsg.Payload())
	authorization.Signatures = []DWNJWSSig{}
	protectedHeaders, _ := jwsMsg.Signatures()[0].ProtectedHeaders().MarshalJSON()

	authorization.Signatures = append(authorization.Signatures, DWNJWSSig{
		Signature: base64.RawURLEncoding.EncodeToString(jwsMsg.Signatures()[0].Signature()),
		Protected: base64.RawURLEncoding.EncodeToString(protectedHeaders),
	})

	return authorization

}
