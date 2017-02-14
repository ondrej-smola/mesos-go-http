package mesos_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("ResourcesBench", func() {
	Measure("BenchmarkPrecisionScalarMath", func(b Benchmarker) {
		count := 10000
		runtime := b.Time("runtime", func() {
			var (
				start   = resources(resource(name("cpus"), valueScalar(1.001)))
				current = start.Clone()
			)
			for i := 0; i < count; i++ {
				current = current.Plus(current...).Plus(current...).Minus(current...).Minus(current...)
			}
		})
		Expect(runtime.Seconds()).
			Should(BeNumerically("<", time.Duration(count)*20*time.Microsecond), "Should take at most 20 micros per operation")
	}, 10)
})
