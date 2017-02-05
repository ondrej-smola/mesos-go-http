package recordio_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRecordio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Recordio Suite")
}
