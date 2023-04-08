package client

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	crypto2 "github.com/TBD54566975/ssi-sdk/crypto"
)

type Keypair struct {
	KeyType    crypto2.KeyType
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

func NewKeypair() Keypair {

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	return Keypair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

}

func FromEd25519PrivateKey(privateKey ed25519.PrivateKey) Keypair {

	var pkCrypto crypto.PrivateKey = privateKey

	return Keypair{
		PrivateKey: pkCrypto,
		PublicKey:  privateKey.Public().(ed25519.PublicKey),
	}

}
