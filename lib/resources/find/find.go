package find

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/ondrej-smola/mesos-go-http/lib"
	"github.com/ondrej-smola/mesos-go-http/lib/resources/filter"
)

// overwrite only in init()
var RandFunc = func(i int) int {
	return rand.Intn(i)
}

func Scalar(name mesos.ResourceName, value float64, in ...*mesos.Resource) (*mesos.Resource, mesos.Resources, bool) {
	res := mesos.Resources(in).Clone()

	cpus, remaining := filter.First(filter.MinScalarValue(name, value), res...)
	if cpus == nil {
		return nil, remaining, false
	}

	if cpus.ScalarValueOrZero() > value {
		remCpus := cpus.Clone()
		remCpus.Scalar.Value = mesos.F64p(mesos.SubtractFixed(cpus.ScalarValueOrZero(), value))
		// add back cpu surplus
		remaining = append(remaining, remCpus)

		cpus.Scalar.Value = mesos.F64p(value)
	}

	return cpus, remaining, true
}

// select first element from random range in provided resources
func RandomInRange(name mesos.ResourceName, in ...*mesos.Resource) (*mesos.Resource, mesos.Resources, bool) {
	allResources := mesos.Resources(in).Clone()
	f := filter.And(filter.Name(name), filter.Range())

	ranges := filter.All(f, allResources...)
	others := filter.All(filter.Not(f), allResources...)

	for i, r := range ranges {
		resRanges := r.RangesOrZero()
		if len(resRanges) == 0 {
			others = append(others, ranges[i])
			continue
		}

		// select random range
		idx := RandFunc(len(resRanges))

		selectedRange := resRanges[idx]
		var rem *mesos.Resource

		if selectedRange.Len() == 1 {
			// selected range is result

			// create copy with remaining free ranges
			rem = r.Clone()
			rem.Ranges = &mesos.Value_Ranges{Range: append(resRanges[:idx], resRanges[idx+1:]...)}
			// set current range to be result with selected range
			r.Ranges = &mesos.Value_Ranges{Range: []*mesos.Value_Range{selectedRange}}

		} else {
			// first element from range is result

			// create remaining ranges
			rCpy := selectedRange.Clone()
			rCpy.Begin = mesos.UI64p(selectedRange.GetBegin() + 1)

			remRanges := append(resRanges[:idx], rCpy)
			remRanges = append(remRanges, resRanges[idx+1:]...)

			// create resource copy with remaining ranges
			rem = r.Clone()
			rem.Ranges = &mesos.Value_Ranges{Range: remRanges}

			// set selected range to contain only selected element
			selectedRange.End = mesos.UI64p(selectedRange.GetBegin())
		}

		others = append(others, rem)
		others = append(others, ranges[i+1:]...)

		r.Ranges = &mesos.Value_Ranges{Range: []*mesos.Value_Range{selectedRange}}

		return r, others, true
	}

	return nil, mesos.Resources(in), false
}

func ValuesInRange(name mesos.ResourceName, values []uint64, in ...*mesos.Resource) (mesos.Resources, mesos.Resources, bool) {
	res := mesos.Resources(in).Clone()

	f := filter.And(filter.Name(name), filter.Range())

	ranges := filter.All(f, res...)
	rem := filter.All(filter.Not(f), res...)

	// sort values
	toFind := append([]uint64(nil), values...)
	sort.Slice(toFind, func(i, j int) bool {
		return toFind[i] < toFind[j]
	})

	takeRanges := mesos.Resources{}

	for i, res := range ranges {
		take := mesos.Ranges{}
		skip := mesos.Ranges{}

		toProcess := res.RangesOrZero()
		toProcess.Sort()

		for len(toProcess) > 0 {
			if len(toFind) == 0 {
				fmt.Println("append", toProcess)
				skip = append(skip, toProcess...)
				break
			}

			// pop next
			currentRange := toProcess[0]
			toProcess = toProcess[1:]

			for _, v := range toFind {
				// split
				l, r, ok := currentRange.Split(v)
				if ok {
					// prepend right range to process list
					if r != nil {
						toProcess = append(mesos.Ranges{r}, toProcess...)
					}
					// left range is appended to skip list
					if l != nil {
						skip = append(skip, l)
					}
					// create value range with single found element
					take = append(take, &mesos.Value_Range{
						Begin: mesos.UI64p(v),
						End:   mesos.UI64p(v),
					})
					toFind = toFind[1:]
					break
				} else {
					skip = append(skip, currentRange)
				}
			}
		}

		if len(take) > 0 {
			// create clone with ranges to use
			cl := res.Clone()
			cl.Ranges = &mesos.Value_Ranges{Range: take}
			takeRanges = append(takeRanges, cl)
		}

		if len(skip) > 0 {
			// create clone with ranges to skip
			cl := res.Clone()
			cl.Ranges = &mesos.Value_Ranges{Range: skip}
			rem = append(rem, cl)
		}

		if len(toFind) == 0 {
			// append all other
			rem = append(rem, ranges[i+1:]...)
			break
		}
	}

	if len(toFind) == 0 {
		return takeRanges, rem, true
	} else {
		return nil, mesos.Resources(in), false
	}
}
