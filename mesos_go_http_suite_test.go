package mesos_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMesosGoHttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MesosGoHttp Suite")
}
