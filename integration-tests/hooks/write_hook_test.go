package hooks_test

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openreserveio/dwn/go/client"
	"os"
)

var _ = Describe("WriteHook", func() {

	DWN_HOST := os.Getenv("DWN_API_HOST")
	DWN_PORT := os.Getenv("DWN_API_PORT")
	TEST_SCHEMA := "https://openreserve.io/schemas/test-hooks.json"

	Describe("SaveHook", func() {

		Context("When a record and its hook is saved", func() {

			dwnClient := client.CreateDWNClient(fmt.Sprintf("http://%s:%s/", DWN_HOST, DWN_PORT))

			authorKeypair := client.New()
			authorIdentity := client.FromKeypair(authorKeypair)

			recipientKeypair := client.New()
			recipientIdentity := client.FromKeypair(recipientKeypair)

			bodyString := fmt.Sprintf("{\"name\":\"testing things\", \"random\":\"%s\"}", uuid.NewString())
			body := []byte(bodyString)

			var recordId string
			var err error

			It("Should save the record", func() {
				println(fmt.Sprintf("Record Body:  %s", bodyString))
				recordId, err = dwnClient.SaveData(TEST_SCHEMA, body, client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &authorIdentity, &recipientIdentity)
				Expect(err).To(BeNil())
				Expect(recordId).ToNot(BeEmpty())
				println(fmt.Sprintf("Record ID:  %s", recordId))
			})

			It("Should save the hook for the record", func() {
				println(fmt.Sprintf("Record ID for HOOK:  %s", recordId))
				hookId, err := dwnClient.SaveHook(TEST_SCHEMA, recordId, "http://localhost:8080", &authorIdentity)
				Expect(err).To(BeNil())
				Expect(hookId).ToNot(BeEmpty())
				println(fmt.Sprintf("Hook ID:  %s", hookId))
			})

			It("Should create a callback if we update the record", func() {

				err := dwnClient.UpdateData(TEST_SCHEMA, recordId, []byte("YAAAAAAAAS"), client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &authorIdentity)
				Expect(err).To(BeNil())

			})

		})

	})

})
