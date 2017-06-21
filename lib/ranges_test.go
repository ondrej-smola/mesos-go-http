package mesos_test

import (
	. "github.com/ondrej-smola/mesos-go-http/lib"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ranges", func() {

	It("Split", func() {

		var tests = []struct {
			in    *Value_Range
			at    uint64
			ok    bool
			left  *Value_Range
			right *Value_Range
		}{
			{NewRange(0, 10), 5, true, NewRange(0, 4), NewRange(6, 10)},
			{NewRange(5, 10), 0, false, nil, nil},
			{NewRange(5, 10), 15, false, nil, nil},
			{NewRange(5, 5), 5, true, nil, nil},
			{NewRange(5, 9), 9, true, NewRange(5, 8), nil},
			{NewRange(5, 9), 5, true, nil, NewRange(6, 9)},
		}

		for i, tt := range tests {
			l, r, ok := tt.in.Split(tt.at)
			Expect(ok).To(Equal(tt.ok),
				fmt.Sprintf("[%v] Split: %v at %v: should be ok '%v' (expected '%v')",
					i, tt.in, tt.at, ok, tt.ok))

			Expect(l).To(Equal(tt.left),
				fmt.Sprintf("[%v] Split: %v at %v: left is '%v' (expected '%v')",
					i, tt.in, tt.at, l, tt.left))

			Expect(r).To(Equal(tt.right),
				fmt.Sprintf("[%v] Split: %v at %v: right is '%v' (expected '%v')",
					i, tt.in, tt.at, r, tt.right))
		}
	})

	It("Contains", func() {
		var tests = []struct {
			in   *Value_Range
			what uint64
			ok   bool
		}{
			{NewRange(0, 10), 5, true},
			{NewRange(5, 10), 0, false},
			{NewRange(5, 10), 15, false},
			{NewRange(5, 5), 5, true},
			{NewRange(5, 9), 9, true},
			{NewRange(5, 9), 5, true},
		}

		for i, tt := range tests {
			ok := tt.in.Contains(tt.what)
			Expect(ok).To(Equal(tt.ok),
				fmt.Sprintf("[%v] Contains: %v what %v: should be ok '%v' (expected '%v')",
					i, tt.in, tt.what, ok, tt.ok))
		}
	})
})
