package find_test

import (
	. "github.com/ondrej-smola/mesos-go-http/lib/resources"
	. "github.com/ondrej-smola/mesos-go-http/lib/resources/find"

	"fmt"

	"github.com/ondrej-smola/mesos-go-http/lib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Find", func() {

	// override rand function to always return last element
	RandFunc = func(i int) int {
		return i - 1
	}

	It("Scalar", func() {
		var tests = []struct {
			name       mesos.ResourceName
			in         mesos.Resources
			value      float64
			shouldFind bool
			find       *mesos.Resource
			rem        mesos.Resources
		}{
			{
				mesos.CPUS,
				mesos.Resources{Cpus(1).WithRole("my_role"), Cpus(1), Mem(256)},
				0.5,
				true,
				Cpus(0.5).WithRole("my_role"),
				mesos.Resources{Cpus(0.5).WithRole("my_role"), Cpus(1), Mem(256)},
			},
			{
				mesos.CPUS,
				mesos.Resources{Cpus(1), Cpus(1), Mem(256)},
				1,
				true,
				Cpus(1),
				mesos.Resources{Cpus(1), Mem(256)},
			},
			{
				mesos.MEM,
				mesos.Resources{Cpus(1), Cpus(1), Mem(128), Mem(256)},
				512,
				false,
				nil,
				mesos.Resources{Cpus(1), Cpus(1), Mem(128), Mem(256)},
			},
		}

		for i, tt := range tests {
			found, rem, ok := Scalar(tt.name, tt.value, tt.in...)
			Expect(ok).To(Equal(tt.shouldFind),
				fmt.Sprintf("[%v] Find[%v]: in %v, value %v: should found '%v' (expected '%v')",
					i, tt.name, tt.in, tt.value, ok, tt.shouldFind))

			Expect(found).To(Equal(tt.find),
				fmt.Sprintf("[%v] Find[%v]: in %v, value %v: found %v  (expected %v)",
					i, tt.name, tt.in, tt.value, found, tt.find))

			Expect(rem).To(ConsistOf(tt.rem),
				fmt.Sprintf("[%v] Find[%v]: in %v, value %v: rem %v (expected %v)",
					i, tt.name, tt.in, tt.value, rem, tt.rem))

		}
	})

	It("RandomInRange", func() {
		var tests = []struct {
			name       mesos.ResourceName
			in         mesos.Resources
			shouldFind bool
			find       *mesos.Resource
			rem        mesos.Resources
		}{
			{
				mesos.PORTS,
				mesos.Resources{Ports(mesos.NewRange(2, 3)).WithRole("my_role"), Cpus(1), Mem(256)},
				true,
				Ports(mesos.NewRange(2, 2)).WithRole("my_role"),
				mesos.Resources{Ports(mesos.NewRange(3, 3)).WithRole("my_role"), Cpus(1), Mem(256)},
			},
			{
				mesos.PORTS,
				mesos.Resources{
					Ports(mesos.NewRange(2, 3), mesos.NewRange(4, 4), mesos.NewRange(5, 5)).WithRole("my_role"),
					Ports(mesos.NewRange(5, 6)),
					Cpus(1),
					Mem(256)},
				true,
				Ports(mesos.NewRange(5, 5)).WithRole("my_role"),
				mesos.Resources{
					Ports(mesos.NewRange(2, 3), mesos.NewRange(4, 4)).WithRole("my_role"),
					Ports(mesos.NewRange(5, 6)),
					Cpus(1),
					Mem(256),
				},
			},
			{
				mesos.PORTS,
				mesos.Resources{Cpus(1), Mem(256)},
				false,
				nil,
				mesos.Resources{Cpus(1), Mem(256)},
			},
		}

		for i, tt := range tests {
			found, rem, ok := RandomInRange(tt.name, tt.in...)
			Expect(ok).To(Equal(tt.shouldFind),
				fmt.Sprintf("[%v] RandomInRange[%v]: in %v: should found '%v' (expected '%v')",
					i, tt.name, tt.in, ok, tt.shouldFind))

			Expect(tt.find).To(Equal(found),
				fmt.Sprintf("[%v] RandomInRange[%v]: in %v: found %v  (expected  %v)",
					i, tt.name, tt.in, found, tt.find))

			Expect(rem).To(ConsistOf(tt.rem),
				fmt.Sprintf("[%v] RandomInRange[%v]: in %v, rem %v (expected %v)",
					i, tt.name, tt.in, rem, tt.rem))
		}
	})

	It("ValuesInRange", func() {
		var tests = []struct {
			name       mesos.ResourceName
			in         mesos.Resources
			values     []uint64
			shouldFind bool
			find       mesos.Resources
			rem        mesos.Resources
		}{
			{
				mesos.PORTS,
				mesos.Resources{Ports(mesos.NewRange(2, 9), mesos.NewRange(15, 20)).WithRole("my_role"), Cpus(1)},
				[]uint64{2, 7, 20},
				true,
				mesos.Resources{Ports(mesos.NewRange(2, 2), mesos.NewRange(7, 7), mesos.NewRange(20, 20)).WithRole("my_role")},
				mesos.Resources{Ports(mesos.NewRange(3, 6), mesos.NewRange(8, 9), mesos.NewRange(15, 19)).WithRole("my_role"), Cpus(1)},
			},
			{
				mesos.PORTS,
				mesos.Resources{Ports(mesos.NewRange(2, 2), mesos.NewRange(15, 15))},
				[]uint64{2, 15},
				true,
				mesos.Resources{Ports(mesos.NewRange(2, 2), mesos.NewRange(15, 15))},
				mesos.Resources{},
			},
			{
				mesos.PORTS,
				mesos.Resources{
					Ports(mesos.NewRange(2, 2), mesos.NewRange(14, 16)),
					Ports(mesos.NewRange(18, 20)).WithRole("my_role")},
				[]uint64{2, 16, 20},
				true,
				mesos.Resources{
					Ports(mesos.NewRange(2, 2), mesos.NewRange(16, 16)),
					Ports(mesos.NewRange(20, 20)).WithRole("my_role"),
				},
				mesos.Resources{
					Ports(mesos.NewRange(14, 15)),
					Ports(mesos.NewRange(18, 19)).WithRole("my_role"),
				},
			},
			{
				mesos.PORTS,
				mesos.Resources{Cpus(1)},
				[]uint64{1},
				false,
				nil,
				mesos.Resources{Cpus(1)},
			},
			{
				mesos.PORTS,
				mesos.Resources{Cpus(1)},
				nil,
				true,
				mesos.Resources{},
				mesos.Resources{Cpus(1)},
			},
		}

		for i, tt := range tests {
			found, rem, ok := ValuesInRange(tt.name, tt.values, tt.in...)
			Expect(ok).To(Equal(tt.shouldFind),
				fmt.Sprintf("[%v] ValuesInRange[%v]: in %v: values %v: should found '%v' (expected '%v')",
					i, tt.name, tt.in, tt.values, ok, tt.shouldFind))

			Expect(tt.find).To(Equal(found),
				fmt.Sprintf("[%v] ValuesInRange[%v]: in %v: values %v: found %v  (expected  %v)",
					i, tt.name, tt.in, tt.values, found, tt.find))

			Expect(rem).To(ConsistOf(tt.rem),
				fmt.Sprintf("[%v] ValuesInRange[%v]: in %v, values %v: rem %v (expected %v)",
					i, tt.name, tt.in, tt.values, rem, tt.rem))
		}

	})
})
