package collection_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/applications/dwn/service/collsvc/collection"
	"github.com/openreserveio/dwn/go/did"
	"github.com/openreserveio/dwn/go/generated/mocks"
	"github.com/openreserveio/dwn/go/model"
	"time"
)

var _ = Describe("StoreCollection", func() {

	mockController := gomock.NewController(GinkgoT())

	Context("Storing a new record (without a parent)", func() {

		collectionStore := mocks.NewMockCollectionStore(mockController)

		It("Should have tried to store the record with an Initial Entry", func() {

			authorPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
			authorPublicKey := authorPrivateKey.PublicKey
			authorDID, _ := did.CreateKeyDID(&authorPublicKey)

			recipientPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
			recipientPublicKey := recipientPrivateKey.PublicKey
			recipientDID, _ := did.CreateKeyDID(&recipientPublicKey)

			body := []byte("{\"name\":\"test\"}")
			message := model.CreateMessage(authorDID, recipientDID, model.DATA_FORMAT_JSON, body, model.METHOD_COLLECTIONS_WRITE, "", "")
			message.Descriptor.DateCreated = time.Now()

			descriptorCID := model.CreateDescriptorCID(message.Descriptor)
			processingCID := model.CreateProcessingCID(message.Processing)
			recordId := model.CreateRecordCID(descriptorCID, processingCID)
			message.RecordID = recordId

			collectionStore.EXPECT().GetCollectionRecord(recordId)
			collectionStore.EXPECT().CreateCollectionRecord(gomock.Any(), gomock.Any()).Return(nil)

			res, err := collection.StoreCollection(context.Background(), collectionStore, message)
			Expect(err).To(BeNil())
			Expect(res).ToNot(BeNil())
			Expect(res.Status).To(Equal("OK"))

		})

	})
	//
	//Context("Saving a collection with an existing ID", func() {
	//
	//	collectionStore := mocks.NewMockCollectionStore(mockController)
	//	collectionStore.EXPECT().PutCollectionItem(gomock.Any()).Return(nil)
	//
	//	It("Should have saved the collection item", func() {
	//
	//		collectionItemId := primitive.NewObjectID().Hex()
	//		newId, ownerDID, err := collection.StoreCollection(collectionStore, "https://openreserve.io/schemas/test.json", collectionItemId, "", "", "", []byte("1"), "did:test:test1", "did:test:test2")
	//		Expect(err).To(BeNil())
	//		Expect(newId).ToNot(BeEmpty())
	//		Expect(newId).To(Equal(collectionItemId))
	//		collectionStore.EXPECT().PutCollectionItem(gomock.AssignableToTypeOf(storage.CollectionItem{})).AnyTimes()
	//		Expect(ownerDID).To(Equal("did:test:test2"))
	//	})
	//
	//})

})
