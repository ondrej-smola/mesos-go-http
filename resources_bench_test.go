package mesos_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		// should take at most 20 micros per operation
		Expect(runtime.Seconds()).Should(BeNumerically("<", 0.2), "SomethingHard() shouldn't take too long.")
	}, 10)
})
