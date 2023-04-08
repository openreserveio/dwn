package client

import (
	"github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did"
)

type Identity struct {
	DID     string
	DIDKey  *did.DIDKey
	Keypair Keypair
}

func FromKeypair(keypair Keypair) Identity {

	didKey, err := did.CreateDIDKey(crypto.Ed25519, keypair.PublicKey.([]byte))
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
