package model

import (
	"context"
	didsdk "github.com/TBD54566975/ssi-sdk/did"
	"strings"
)

var didResolver *didsdk.MultiMethodResolver

func init() {

	// Resolve the signing key
	resolvers := []didsdk.Resolver{
		didsdk.KeyResolver{},
		didsdk.WebResolver{},
		didsdk.PKHResolver{},
		didsdk.PeerResolver{},
	}
	resolver, err := didsdk.NewResolver(resolvers...)
	if err != nil {
		panic(err)
	}

	didResolver = resolver

}

func ResolveDID(didToResolve string) (*didsdk.Document, error) {

	// Strip out everything after a # if it exists
	if i := strings.Index(didToResolve, "#"); i != -1 {
		didToResolve = didToResolve[:i]
	}

	res, err := didResolver.Resolve(context.Background(), didToResolve)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	return &res.Document, nil

}
