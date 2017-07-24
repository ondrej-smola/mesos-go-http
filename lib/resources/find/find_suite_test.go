package find_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFind(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Find Suite")
}
