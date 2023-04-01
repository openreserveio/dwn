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
	TEST_PROTOCOL := "https://openreserve.io/protocol/test-protocol.json"
	TEST_PROTOCOL_VERSION := "0.0.1"

	Describe("SaveHookForRecord", func() {

		Context("When a record and its hook is saved", func() {

			dwnClient := client.CreateDWNClient(fmt.Sprintf("http://%s:%s/", DWN_HOST, DWN_PORT))

			authorKeypair := client.NewKeypair()
			authorIdentity := client.FromKeypair(authorKeypair)

			recipientKeypair := client.NewKeypair()
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
				hookId, err := dwnClient.SaveHookForRecord(TEST_SCHEMA, recordId, "http://localhost:8080", &authorIdentity)
				Expect(err).To(BeNil())
				Expect(hookId).ToNot(BeEmpty())
				println(fmt.Sprintf("Hook ID:  %s", hookId))
			})

			It("Should create a callback if we update the record", func() {

				_, err := dwnClient.UpdateData(TEST_SCHEMA, recordId, recordId, []byte("YAAAAAAAAS"), client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &authorIdentity)
				Expect(err).To(BeNil())

			})

		})

	})

	Describe("SaveHookForSchemaAndProtocol", func() {

		Context("A hook for a schema and protocol are defined, then a record for that schema and protcol is created", func() {

			dwnClient := client.CreateDWNClientForProtocol(fmt.Sprintf("http://%s:%s/", DWN_HOST, DWN_PORT), TEST_PROTOCOL, TEST_PROTOCOL_VERSION)

			authorKeypair := client.NewKeypair()
			authorIdentity := client.FromKeypair(authorKeypair)

			recipientKeypair := client.NewKeypair()
			recipientIdentity := client.FromKeypair(recipientKeypair)

			bodyString := fmt.Sprintf("{\"name\":\"testing things\", \"random\":\"%s\"}", uuid.NewString())
			body := []byte(bodyString)

			var recordId string
			var err error

			It("Should save the hook for the record", func() {
				println(fmt.Sprintf("Record ID for HOOK:  %s", recordId))
				hookId, err := dwnClient.SaveHookForSchemaAndProtocol(TEST_SCHEMA, TEST_PROTOCOL, TEST_PROTOCOL_VERSION, "http://localhost:8080", &authorIdentity)
				Expect(err).To(BeNil())
				Expect(hookId).ToNot(BeEmpty())
				println(fmt.Sprintf("Hook ID:  %s", hookId))
			})

			It("Should save a record with that schema, protocol, and protocol version", func() {
				println(fmt.Sprintf("Record Body:  %s", bodyString))
				recordId, err = dwnClient.SaveData(TEST_SCHEMA, body, client.HEADER_CONTENT_TYPE_APPLICATION_JSON, &authorIdentity, &recipientIdentity)
				Expect(err).To(BeNil())
				Expect(recordId).ToNot(BeEmpty())
				println(fmt.Sprintf("Record ID:  %s", recordId))
			})

		})

	})

})
