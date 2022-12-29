package did

import (
	"github.com/openreserveio/dwn/go/log"
	"strings"
)

func ResolvePublicKey(didString string) any {

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
	case DID_METHOD_KEY:
		return ParseKeyFromKeyDID(didParts[2])

	case DID_METHOD_WEB:
		log.Error("Support for did:web coming soon!")
		return nil

	default:
		log.Error("Only support for did:key")
		return nil
	}

}
