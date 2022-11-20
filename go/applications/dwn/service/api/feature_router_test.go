package api_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/applications/dwn/service/api"
	"github.com/openreserveio/dwn/go/generated/mocks"
	"github.com/openreserveio/dwn/go/model"
)

var _ = Describe("Feature Router", func() {

	mockController := gomock.NewController(GinkgoT())

	Context("Simple Routing/Reply", func() {

		var err error
		var router *api.FeatureRouter
		mockCollSvcClient := mocks.NewMockCollectionServiceClient(mockController)

		It("Should create a feature router instance", func() {
			router, err = api.CreateFeatureRouter(mockCollSvcClient, 15)
			Expect(err).To(BeNil())
			Expect(router).ToNot(BeNil())
		})

		It("Should route single messages", func() {

			ro := model.RequestObject{
				Messages: []model.Message{
					model.Message{
						RecordID: "TEST",
						Data:     "TEST",
						Processing: model.MessageProcessing{
							Nonce:        "TEST",
							AuthorDID:    "did:test:test",
							RecipientDID: "did:test:test",
						},
						Descriptor: model.Descriptor{
							Nonce:      "TEST",
							Method:     "TEST",
							DataCID:    "TEST",
							DataFormat: "TEST",
						},
					},
				},
			}

			resp, err := router.Route(&ro)

			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(resp.(model.ResponseObject).Status.Code).To(Equal(200))

		})

	})

})
