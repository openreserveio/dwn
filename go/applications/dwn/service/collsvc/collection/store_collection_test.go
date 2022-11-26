package collection_test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("StoreCollection", func() {

	// mockController := gomock.NewController(GinkgoT())

	//Context("Storing a new collection (without an ID)", func() {
	//
	//	collectionStore := mocks.NewMockCollectionStore(mockController)
	//	collectionStore.EXPECT().PutCollectionItem(gomock.Any()).Return(nil)
	//
	//	It("Should have tried to store the collection item with a new ID", func() {
	//
	//		newId, ownerDID, err := collection.StoreCollection(collectionStore, "https://openreserve.io/schemas/test.json", "", "", "", "", []byte("1"), "did:test:test1", "did:test:test2")
	//		Expect(err).To(BeNil())
	//		Expect(newId).ToNot(BeEmpty())
	//		collectionStore.EXPECT().PutCollectionItem(gomock.AssignableToTypeOf(storage.CollectionItem{})).AnyTimes()
	//		Expect(ownerDID).To(Equal("did:test:test2"))
	//
	//	})
	//
	//})
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
