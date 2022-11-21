package collections_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	. "github.com/onsi/ginkgo/v2"
	"github.com/openreserveio/dwn/go/did"
)

var _ = Describe("Write Collection", func() {

	Describe("Writing a brand new collection item", func() {

		authorPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		authorPublicKey := authorPrivateKey.PublicKey
		authorDID := did.CreateKeyDID(&authorPublicKey)

		recipientPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		recipientPublicKey := recipientPrivateKey.PublicKey
		recipientDID := did.CreateKeyDID(&recipientPublicKey)

	})

})
