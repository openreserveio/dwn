package client

import "github.com/openreserveio/dwn/go/did"

type Identity struct {
	DID     string
	Keypair Keypair
}

func FromKeypair(keypair Keypair) Identity {

	ident := Identity{Keypair: keypair}
	identDID, _ := did.CreateKeyDID(keypair.PublicKey)
	ident.DID = identDID

	return ident

}
