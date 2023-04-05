package client

import (
	"crypto/ecdsa"
	"github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did"
)

type Identity struct {
	DID     string
	Keypair Keypair
}

func FromKeypair(keypair Keypair) Identity {

	ident := Identity{Keypair: keypair}
	privateKey, didKey, _ := did.GenerateDIDKey(crypto.Ed25519)

	identity := Identity{
		DID: didKey.String(),
		Keypair: ed2
	}

	return ident

}
