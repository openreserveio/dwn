package model

import (
	"crypto"
	"encoding/base64"
	"fmt"
	didsdk "github.com/TBD54566975/ssi-sdk/did"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"strings"
)

func VerifyAttestation(message *Message) bool {

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

func CreateAttestation(message *Message, authKeyUri string, publicKey crypto.PublicKey, privateKey crypto.PrivateKey) DWNJWS {

	attestation := DWNJWS{}

	// Using github.com/lestrrat-go/jwx package here because *amazing*
	/*
			From Spec:
			The Message object MUST contain an attestation property, and its value MUST be a General object representation of a [RFC7515] JSON Web Signature composed as follows:
		The object MUST include a payload property, and its value MUST be the stringified Version 1 CID of the DAG CBOR encoded descriptor object, whose composition is defined in the Message Descriptor section of this specification.
		The object MUST include a protected property, and its value MUST be an object composed of the following values:
		The object MUST include an alg property, and its value MUST be the string representing the algorithm used to verify the signature (as defined by the [RFC7515] JSON Web Signature specification).
		The object MUST include a kid property, and its value MUST be a DID URL string identifying the key to be used in verifying the signature.
		The object MUST include a signature property, and its value MUST be a signature string produced by signing the protected and payload values, in accordance with the [RFC7515] JSON Web Signature specification.

	*/
	descriptorCid := CreateDescriptorCID(message.Descriptor)
	// TODO:  Include permissionsGrants and attestation logic per spec
	// attestationCid := CreateAttestationCID(message.Attestation)
	// permissionsGrantCid := CreatePermissionsGrantCID(message.PermissionsGrant)

	additionalHeaders := jws.NewHeaders()
	additionalHeaders.Set("kid", authKeyUri)

	signature, err := jws.Sign([]byte(descriptorCid), jwa.EdDSA, privateKey, jws.WithHeaders(additionalHeaders))
	if err != nil {
		panic(err)
	}

	jwsMsg, err := jws.Parse(signature)
	if err != nil {
		panic(err)
	}

	attestation.Payload = descriptorCid
	attestation.Signatures = []DWNJWSSig{}
	protectedHeaders, _ := jwsMsg.Signatures()[0].ProtectedHeaders().MarshalJSON()

	attestation.Signatures = append(attestation.Signatures, DWNJWSSig{
		Signature: base64.RawURLEncoding.EncodeToString(jwsMsg.Signatures()[0].Signature()),
		Protected: base64.RawURLEncoding.EncodeToString(protectedHeaders),
	})

	return attestation

}
