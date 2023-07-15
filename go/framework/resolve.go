package framework

import (
	"context"
	"errors"
	"fmt"
	didsdk "github.com/TBD54566975/ssi-sdk/did"
	"github.com/openreserveio/dwn/go/observability"
)

var didResolver *didsdk.MultiMethodResolver

func init() {

	resolver, err := BuildResolver([]string{"web", "key", "pkh", "peer"})
	if err != nil {
		panic(err)
	}

	didResolver = resolver
}

func ResolveDID(ctx context.Context, didString string) (*didsdk.Document, error) {

	// Instrumentation
	ctx, sp := observability.Tracer().Start(ctx, "framework.ResolveDID")
	defer sp.End()

	sp.AddEvent(fmt.Sprintf("Resolving DID: %s", didString))
	res, err := didResolver.Resolve(ctx, didString)
	if err != nil {
		sp.RecordError(err)
		return nil, err
	}

	if res == nil {
		sp.AddEvent(fmt.Sprintf("DID %s did not resolve.", didString))
		return nil, nil
	}

	if res.Error != nil {

		if res.Error.NotFound {
			sp.AddEvent(fmt.Sprintf("DID %s was not found", didString))
		} else if res.Error.InvalidDID {
			sp.AddEvent(fmt.Sprintf("DID %s was invalid", didString))
		} else if res.Error.RepresentationNotSupported {
			sp.AddEvent(fmt.Sprintf("DID %s representation is not supported", didString))
		} else {
			sp.AddEvent(fmt.Sprintf("DID %s resolution failed: %v", didString, res.Error))
		}

		return nil, nil

	}

	return &res.Document, nil

}

// BuildResolver builds a DID resolver from a list of methods to support resolution for
func BuildResolver(methods []string) (*didsdk.MultiMethodResolver, error) {

	if len(methods) == 0 {
		return nil, errors.New("no methods provided")
	}
	var resolvers []didsdk.Resolver
	for _, method := range methods {
		resolver, err := getKnownResolver(method)
		if err != nil {
			return nil, err
		}
		resolvers = append(resolvers, resolver)
	}
	if len(resolvers) == 0 {
		return nil, errors.New("no resolvers created")
	}
	return didsdk.NewResolver(resolvers...)
}

// all possible resolvers for the DID service
func getKnownResolver(method string) (didsdk.Resolver, error) {
	switch didsdk.Method(method) {
	case didsdk.KeyMethod:
		return new(didsdk.KeyResolver), nil
	case didsdk.WebMethod:
		return new(didsdk.WebResolver), nil
	case didsdk.PKHMethod:
		return new(didsdk.PKHResolver), nil
	case didsdk.PeerMethod:
		return new(didsdk.PeerResolver), nil
	}
	return nil, fmt.Errorf("unsupported method: %s", method)
}
