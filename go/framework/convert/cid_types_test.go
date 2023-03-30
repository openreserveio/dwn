package convert_test

import (
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openreserveio/dwn/go/framework/convert"
)

var _ = Describe("ToUintConverter", func() {

	Context("when input data is a valid uint64 value", func() {

		It("should return a datamodel.Node with the correct value", func() {

			data := uint64(42)
			node, err := convert.ToUintConverter(data)

			Expect(err).To(BeNil())
			Expect(node.Kind()).To(Equal(datamodel.Kind_Int))

			value, err := node.AsInt()
			Expect(err).To(BeNil())
			Expect(value).To(Equal(int64(data)))

			// Create the expected node using basicnode and compare the values.
			expectedNodeBuilder := basicnode.Prototype.Int.NewBuilder()
			err = expectedNodeBuilder.AssignInt(int64(data))
			Expect(err).To(BeNil())
			expectedNode := expectedNodeBuilder.Build()

			Expect(datamodel.DeepEqual(node, expectedNode)).To(BeTrue())
		})
	})

	Context("when input data is not a uint64 value", func() {
		It("should return an error", func() {
			data := "not a uint64"
			node, err := convert.ToUintConverter(data)

			Expect(err).To(HaveOccurred())
			Expect(node).To(BeNil())
		})
	})
})
