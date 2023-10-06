package resolving_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/framework"
	"github.com/openreserveio/dwn/go/observability"
)

var _ = Describe("ResolveWeb", func() {

	ctx := context.Background()
	sd, _ := observability.InitProviderWithOTELExporter(ctx, "ResolveWeb Integration Test")
	defer sd(ctx)

	It("Resolves a valid web DID", func() {

		doc, err := framework.ResolveDID(context.Background(), "did:web:openreserve.io")
		Expect(err).To(BeNil())
		Expect(doc).ToNot(BeNil())

		Expect(doc.ID).To(Equal("did:web:openreserve.io"))
		Expect(len(doc.VerificationMethod)).To(Equal(3))
		Expect(doc.VerificationMethod[0].ID).To(Equal("did:web:openreserve.io#key-0"))
		Expect(doc.VerificationMethod[1].ID).To(Equal("did:web:openreserve.io#key-1"))
		Expect(doc.VerificationMethod[2].ID).To(Equal("did:web:openreserve.io#key-2"))

	})

})
