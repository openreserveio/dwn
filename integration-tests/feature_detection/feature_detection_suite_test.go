package feature_detection_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFeatureDetection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FeatureDetection Suite")
}
