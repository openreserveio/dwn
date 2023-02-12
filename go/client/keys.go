package client

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type Keypair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func New() Keypair {

	privateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	publicKey := &privateKey.PublicKey

	return Keypair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

}
