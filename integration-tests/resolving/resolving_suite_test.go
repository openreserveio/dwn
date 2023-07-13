package resolving_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestResolving(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resolving Suite")
}
