package backoff_test

import (
	"context"
	"testing"
	"time"

	. "github.com/ondrej-smola/mesos-go-http/backoff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRetry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Backoff Suite")
}

const tolerance = 2 * time.Millisecond

var _ = Describe("Backoff", func() {
	It("Attempts", func() {
		prov := New(
			WithJitterFraction(0),
			WithMaxAttempts(4),
			WithMinWait(10*time.Millisecond),
			WithMaxWait(50*time.Millisecond),
			WithBackoffFactor(2),
		).New()
		defer prov.Close()

		diffs := make([]time.Duration, 5)

		last := time.Now()
		for a := range prov.Attempts() {
			cur := time.Now()
			diffs[a-1] = cur.Sub(last)
			last = cur
		}

		Expect(diffs[0]).To(BeNumerically("~", 0, tolerance))
		Expect(diffs[1]).To(BeNumerically("~", 10*time.Millisecond, tolerance))
		Expect(diffs[2]).To(BeNumerically("~", 20*time.Millisecond, tolerance))
		Expect(diffs[3]).To(BeNumerically("~", 40*time.Millisecond, tolerance))
	})

	It("Reset", func(done Done) {
		prov := New(
			WithJitterFraction(0),
			WithMaxAttempts(5),
			WithMinWait(10*time.Millisecond),
			WithMaxWait(50*time.Millisecond),
			WithBackoffFactor(1.5),
		).New()
		defer prov.Close()

		diffs := make([]time.Duration, 5)
		atts := []int{}
		last := time.Now()

		for a := range prov.Attempts() {
			atts = append(atts, a)
			cur := time.Now()
			diffs[a-1] = cur.Sub(last)
			if a == 3 && len(atts) == 3 {
				prov.Reset()
			}
			last = time.Now()
		}

		Expect(atts).To(Equal([]int{1, 2, 3, 1, 2, 3, 4, 5}))
		Expect(diffs[0]).To(BeNumerically("~", 0, tolerance))
		Expect(diffs[1]).To(BeNumerically("~", 10*time.Millisecond, tolerance))
		Expect(diffs[2]).To(BeNumerically("~", 15*time.Millisecond, tolerance))
		Expect(diffs[3]).To(BeNumerically("~", 22*time.Millisecond, tolerance))
		Expect(diffs[4]).To(BeNumerically("~", 34*time.Millisecond, tolerance))
		close(done)
	})

	It("Context stop", func(done Done) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		prov := New(
			WithJitterFraction(0),
			WithMaxAttempts(5),
			WithMinWait(10*time.Millisecond),
			WithMaxWait(50*time.Millisecond),
			WithBackoffFactor(2),
		).New(ctx)
		defer prov.Close()

		diffs := []time.Duration{}
		last := time.Now()

		for a := range prov.Attempts() {
			cur := time.Now()
			diffs = append(diffs, cur.Sub(last))
			if a == 3 {
				cancel()
			}
			last = time.Now()
		}

		Expect(diffs).To(HaveLen(3))
		Expect(diffs[0]).To(BeNumerically("~", 0, tolerance))
		Expect(diffs[1]).To(BeNumerically("~", 10*time.Millisecond, tolerance))
		Expect(diffs[2]).To(BeNumerically("~", 20*time.Millisecond, tolerance))
		close(done)
	})

	It("Jitter", func(done Done) {
		prov := New(
			WithJitterFraction(0.1),
			WithMaxAttempts(5),
			WithMinWait(10*time.Millisecond),
			WithMaxWait(50*time.Millisecond),
			WithRandFunc(func() float64 {
				return 1
			}),
			WithBackoffFactor(2),
		).New()
		defer prov.Close()

		diffs := make([]time.Duration, 5)
		last := time.Now()

		for a := range prov.Attempts() {
			cur := time.Now()
			diffs[a-1] = cur.Sub(last)
			last = time.Now()
		}

		Expect(diffs[0]).To(BeNumerically("~", 0, tolerance))
		Expect(diffs[1]).To(BeNumerically("~", 11*time.Millisecond, tolerance))
		Expect(diffs[2]).To(BeNumerically("~", 22*time.Millisecond, tolerance))
		Expect(diffs[3]).To(BeNumerically("~", 44*time.Millisecond, tolerance))
		Expect(diffs[4]).To(BeNumerically("~", 50*time.Millisecond, tolerance))

		close(done)
	})
})
