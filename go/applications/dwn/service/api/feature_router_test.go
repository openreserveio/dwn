package api_test

import (
	"encoding/base64"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/applications/dwn/service/api"
	"github.com/openreserveio/dwn/go/generated/mocks"
	"github.com/openreserveio/dwn/go/model"
	"net/http"
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
			Expect(resp.(*model.ResponseObject).Status.Code).To(Equal(200))

		})

		It("Should route multiple messages", func() {

			ro := model.RequestObject{
				Messages: []model.Message{
					model.Message{
						RecordID: "TEST",
						Data:     "TEST",
						Processing: model.MessageProcessing{
							Nonce:        "TEST",
							AuthorDID:    "did:test:test1",
							RecipientDID: "did:test:test2",
						},
						Descriptor: model.Descriptor{
							Nonce:      "TEST",
							Method:     "TEST",
							DataCID:    "TEST",
							DataFormat: "TEST",
						},
					},
					model.Message{
						RecordID: "TEST",
						Data:     "TEST",
						Processing: model.MessageProcessing{
							Nonce:        "TEST",
							AuthorDID:    "did:test:test3",
							RecipientDID: "did:test:test4",
						},
						Descriptor: model.Descriptor{
							Nonce:      "TEST",
							Method:     "TEST",
							DataCID:    "TEST",
							DataFormat: "TEST",
						},
					},
					model.Message{
						RecordID: "TEST",
						Data:     "TEST",
						Processing: model.MessageProcessing{
							Nonce:        "TEST",
							AuthorDID:    "did:test:test5",
							RecipientDID: "did:test:test6",
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
			respObject := resp.(*model.ResponseObject)
			Expect(respObject.Status.Code).To(Equal(200))
			Expect(len(respObject.Replies)).To(Equal(3))

		})

	})

	Context("Writing to a collection", func() {

		var err error
		var router *api.FeatureRouter
		mockCollSvcClient := mocks.NewMockCollectionServiceClient(mockController)

		It("Should create a feature router instance", func() {
			router, err = api.CreateFeatureRouter(mockCollSvcClient, 15)
			Expect(err).To(BeNil())
			Expect(router).ToNot(BeNil())
		})

		It("Should reject a CollectionsWrite without a SchemaURI", func() {

			ro := model.RequestObject{
				Messages: []model.Message{
					model.Message{
						RecordID: "",
						Data:     base64.URLEncoding.EncodeToString([]byte("{}")),
						Processing: model.MessageProcessing{
							Nonce:        uuid.NewString(),
							AuthorDID:    "did:test:test1",
							RecipientDID: "did:test:test2",
						},
						Descriptor: model.Descriptor{
							Nonce:      uuid.NewString(),
							Method:     "CollectionsWrite",
							DataCID:    "",
							DataFormat: model.DATA_FORMAT_JSON,
							Schema:     "",
						},
					},
				},
			}

			resp, err := router.Route(&ro)
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())

			responseObject := resp.(*model.ResponseObject)
			Expect(responseObject.Status.Code).To(Equal(200))
			Expect(len(responseObject.Replies)).To(Equal(1))
			Expect(responseObject.Replies[0].Status.Code).To(Equal(http.StatusBadRequest))

		})

	})

})
