package testutils

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/multiformats/go-multibase"
)

func CreateKeyDID(publicKey *ecdsa.PublicKey) (string, error) {

	pubkeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	keyBlock := pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: pubkeyBytes,
	}

	pemEncodedPublicKey := pem.EncodeToMemory(&keyBlock)
	keyDidEncodedKey, err := multibase.Encode(multibase.Base64url, pemEncodedPublicKey)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("did:key:%s", keyDidEncodedKey), nil

}
