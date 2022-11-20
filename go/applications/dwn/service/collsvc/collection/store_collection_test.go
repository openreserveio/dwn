package collection_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/applications/dwn/service/collsvc/collection"
	"github.com/openreserveio/dwn/go/generated/mocks"
	"github.com/openreserveio/dwn/go/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("StoreCollection", func() {

	mockController := gomock.NewController(GinkgoT())

	Context("Storing a new collection (without an ID)", func() {

		collectionStore := mocks.NewMockCollectionStore(mockController)
		collectionStore.EXPECT().PutCollectionItem(gomock.Any()).Return(nil)

		It("Should have tried to store the collection item with a new ID", func() {

			newId, err := collection.StoreCollection(collectionStore, "https://openreserve.io/schemas/test.json", "", []byte("1"))
			Expect(err).To(BeNil())
			Expect(newId).ToNot(BeEmpty())
			collectionStore.EXPECT().PutCollectionItem(gomock.AssignableToTypeOf(storage.CollectionItem{})).AnyTimes()

		})

	})

	Context("Saving a collection with an existing ID", func() {

		collectionStore := mocks.NewMockCollectionStore(mockController)
		collectionStore.EXPECT().PutCollectionItem(gomock.Any()).Return(nil)

		It("Should have saved the collection item", func() {

			collectionItemId := primitive.NewObjectID().Hex()
			newId, err := collection.StoreCollection(collectionStore, "https://openreserve.io/schemas/test.json", collectionItemId, []byte("1"))
			Expect(err).To(BeNil())
			Expect(newId).ToNot(BeEmpty())
			Expect(newId).To(Equal(collectionItemId))
			collectionStore.EXPECT().PutCollectionItem(gomock.AssignableToTypeOf(storage.CollectionItem{})).AnyTimes()

		})

	})

})
