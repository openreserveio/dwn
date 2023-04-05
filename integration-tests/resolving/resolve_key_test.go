package resolving_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/client"
	"github.com/openreserveio/dwn/go/framework"
	"github.com/openreserveio/dwn/go/observability"
)

var _ = Describe("ResolveKey", func() {

	ctx := context.Background()
	sd, _ := observability.InitProviderWithJaegerExporter(ctx, "ResolveWeb Integration Test")
	defer sd(ctx)

	kp := client.NewKeypair()
	ident := client.FromKeypair(kp)

	didString := fmt.Sprintf("did:key:%s", ident.DID)

	It("Resolves a valid key DID", func() {

		doc, err := framework.ResolveDID(context.Background(), didString)
		Expect(err).To(BeNil())
		Expect(doc).ToNot(BeNil())

		Expect(doc.ID).To(Equal(didString))
		Expect(len(doc.VerificationMethod)).To(Equal(1))

	})

})
