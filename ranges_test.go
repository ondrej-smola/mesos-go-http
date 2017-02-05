package mesos

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ranges", func() {

	It("NewRanges", func() {
		for _, tt := range []struct {
			ns   []uint64
			want Ranges
		}{
			{[]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, Ranges{{0, 10}}},
			{[]uint64{7, 2, 3, 8, 1, 5, 0, 9, 6, 10, 4}, Ranges{{0, 10}}},
			{[]uint64{0, 0, 1, 1, 2, 2, 3, 3, 4, 4}, Ranges{{0, 4}}},
			{[]uint64{0, 1, 3, 5, 6, 8, 9}, Ranges{{0, 1}, {3, 3}, {5, 6}, {8, 9}}},
			{[]uint64{1}, Ranges{{1, 1}}},
			{[]uint64{}, Ranges{}},
		} {
			Expect(NewRanges(tt.ns...)).To(Equal(tt.want))
		}
	})

	It("NewPortRanges", func() {
		for _, tt := range []struct {
			Ranges
			want Ranges
		}{
			{Ranges{{2, 0}, {3, 10}}, Ranges{{0, 10}}},
			{Ranges{{0, 2}, {3, 10}}, Ranges{{0, 10}}},
			{Ranges{{0, 2}, {1, 10}}, Ranges{{0, 10}}},
			{Ranges{{0, 2}, {4, 10}}, Ranges{{0, 2}, {4, 10}}},
			{Ranges{{10, 0}}, Ranges{{0, 10}}},
			{Ranges{}, Ranges{}},
			{nil, Ranges{}},
		} {
			offer := &Offer{Resources: []Resource{tt.resource("ports")}}
			got := NewPortRanges(offer)
			Expect(got).To(Equal(tt.want))
		}
	})

	It("Size", func() {
		for _, tt := range []struct {
			Ranges
			want uint64
		}{
			{Ranges{}, 0},
			{Ranges{{0, 1}, {2, 10}}, 11},
			{Ranges{{0, 1}, {2, 10}, {11, 100}}, 101},
			{Ranges{{1, 1}, {2, 10}, {11, 100}}, 100},
		} {
			Expect(tt.Size()).To(Equal(tt.want))
		}
	})

	It("Squash", func() {
		for _, tt := range []struct {
			Ranges
			want Ranges
		}{
			{Ranges{}, Ranges{}},
			{Ranges{{0, 1}}, Ranges{{0, 1}}},
			{Ranges{{0, 2}, {1, 5}, {2, 10}}, Ranges{{0, 10}}},
			{Ranges{{0, 2}, {2, 5}, {5, 10}}, Ranges{{0, 10}}},
			{Ranges{{0, 2}, {3, 5}, {6, 10}}, Ranges{{0, 10}}},
			{Ranges{{0, 2}, {4, 11}, {6, 10}}, Ranges{{0, 2}, {4, 11}}},
			{Ranges{{0, 2}, {4, 5}, {6, 7}, {8, 10}}, Ranges{{0, 2}, {4, 10}}},
			{Ranges{{0, 2}, {4, 6}, {8, 10}}, Ranges{{0, 2}, {4, 6}, {8, 10}}},
			{Ranges{{0, 1}, {2, 5}, {4, 8}}, Ranges{{0, 8}}},
		} {
			Expect(tt.Squash()).To(Equal(tt.want))
		}
	})

	It("Search", func() {
		for _, tt := range []struct {
			Ranges
			n    uint64
			want int
		}{
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 0, 0},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 1, 0},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 2, 0},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 3, 1},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 4, 1},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 5, 1},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 6, -1},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 7, 2},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 8, 2},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 9, 2},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 10, 2},
			{Ranges{{0, 2}, {3, 5}, {7, 10}}, 11, -1},
			{Ranges{{0, 2}, {4, 4}, {5, 10}}, 4, 1},
		} {
			Expect(tt.Search(tt.n)).To(Equal(tt.want))
		}
	})

	It("Partition", func() {
		for _, tt := range []struct {
			Ranges
			n     uint64
			want  Ranges
			found bool
		}{
			{Ranges{}, 0, Ranges{}, false},
			{Ranges{{0, 10}, {12, 20}}, 100, Ranges{{0, 10}, {12, 20}}, false},
			{Ranges{{0, 10}, {12, 20}}, 0, Ranges{{1, 10}, {12, 20}}, true},
			{Ranges{{0, 10}, {12, 20}}, 13, Ranges{{0, 10}, {12, 12}, {14, 20}}, true},
			{Ranges{{0, 10}, {12, 20}}, 5, Ranges{{0, 4}, {6, 10}, {12, 20}}, true},
			{Ranges{{0, 10}, {12, 20}}, 19, Ranges{{0, 10}, {12, 18}, {20, 20}}, true},
			{Ranges{{0, 10}, {12, 20}}, 10, Ranges{{0, 9}, {12, 20}}, true},
			{Ranges{{0, 10}, {12, 12}, {14, 20}}, 12, Ranges{{0, 10}, {14, 20}}, true},
		} {
			got, found := tt.Partition(tt.n)
			Expect(got).To(Equal(tt.want))
			Expect(found).To(Equal(tt.found))
		}
	})

	It("MinMax", func() {
		for _, tt := range []struct {
			Ranges
			min, max uint64
		}{
			{Ranges{{1, 10}, {100, 1000}}, 1, 1000},
			{Ranges{{0, 10}, {12, 20}}, 0, 20},
			{Ranges{{5, 10}}, 5, 10},
			{Ranges{{0, 0}}, 0, 0},
		} {
			Expect(tt.Min()).To(Equal(tt.min))
			Expect(tt.Max()).To(Equal(tt.max))
		}
	})
})
