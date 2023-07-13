package client

import (
	"crypto/ed25519"
	"github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did"
)

type Identity struct {
	DID     string
	DIDKey  *did.DIDKey
	Keypair Keypair
}

func FromKeypair(keypair Keypair) Identity {

	var publicKeyBytes []byte
	switch keypair.KeyType {
	case crypto.Ed25519:
		publicKeyBytes = keypair.PublicKey.(ed25519.PublicKey)

	default:
		panic("Unsupported key type")
	}

	didKey, err := did.CreateDIDKey(crypto.Ed25519, publicKeyBytes)
	if err != nil {
		panic(err)
	}

	identity := Identity{
		DID:     didKey.String(),
		DIDKey:  didKey,
		Keypair: keypair,
	}

	return identity

}
