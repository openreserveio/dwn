package collections_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCollections(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collections Suite")
}
