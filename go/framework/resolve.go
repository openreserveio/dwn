package did

import "github.com/TBD54566975/ssi-sdk/did"

type DIDMethodResolver interface {
	Resolve(didString string) did.DIDDocument
}
