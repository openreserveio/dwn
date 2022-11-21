package feature_detection_test

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/model"
)

var _ = Describe("Get Feature Detection", func() {

	client := resty.New()

	It("Gets a list of features available", func() {

		resp, err := client.R().Get("http://localhost:8080/")
		Expect(err).To(BeNil())
		Expect(resp).ToNot(BeNil())

		var fdBody model.FeatureDetection
		err = json.Unmarshal(resp.Body(), &fdBody)
		Expect(err).To(BeNil())
		Expect(fdBody.Type).To(Equal("FeatureDetection"))
		Expect(fdBody.Interfaces.Collections.CollectionsQuery).To(BeTrue())
		Expect(fdBody.Interfaces.Collections.CollectionsWrite).To(BeTrue())
		Expect(fdBody.Interfaces.Collections.CollectionsCommit).To(BeTrue())
		Expect(fdBody.Interfaces.Collections.CollectionsDelete).To(BeTrue())

	})

})
