package did

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/multiformats/go-multibase"
	"github.com/openreserveio/dwn/go/log"
	"strings"
)

func Resolve(didString string) any {

	if didString == "" {
		return nil
	}

	// split the string
	didParts := strings.Split(didString, ":")
	if len(didParts) < 3 || didParts[0] != "did" {
		return nil
	}

	// Get the method
	didMethod := didParts[1]
	switch didMethod {
	case "key":
		multibaseKey := didParts[2]
		_, pemPublicKey, err := multibase.Decode(multibaseKey)
		if err != nil {
			return nil
		}

		publicKey, err := jwt.ParseECPublicKeyFromPEM(pemPublicKey)
		if err != nil {
			return nil
		}
		return publicKey

	default:
		log.Error("Only support for did:key")
		return nil
	}

}
