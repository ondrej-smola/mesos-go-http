package mesos

import (
	"sort"
)

type Ranges []*Value_Range

func NewRange(begin, end uint64) *Value_Range {
	return &Value_Range{
		Begin: &begin,
		End:   &end,
	}
}

func (r *Value_Range) Len() uint64 {
	return r.GetEnd() - r.GetBegin() + 1
}

func (r *Value_Range) Contains(el uint64) bool {
	return el >= r.GetBegin() && el <= r.GetEnd()
}

// Split range around el, does not modify current range
func (r *Value_Range) Split(el uint64) (*Value_Range, *Value_Range, bool) {
	if !r.Contains(el) {
		return nil, nil, false
	}

	begin := r.GetBegin()
	end := r.GetEnd()

	if begin == el && end == el {
		return nil, nil, true
	} else if begin == el {
		return nil, &Value_Range{Begin: UI64p(begin + 1), End: UI64p(end)}, true
	} else if end == el {
		return &Value_Range{Begin: UI64p(begin), End: UI64p(end - 1)}, nil, true
	} else {
		left := &Value_Range{Begin: UI64p(begin), End: UI64p(el - 1)}
		right := &Value_Range{Begin: UI64p(el + 1), End: UI64p(end)}
		return left, right, true
	}
}

func (r Ranges) Sort() {
	sort.Slice(r, func(i, j int) bool {
		return r[i].GetBegin() < r[j].GetBegin() || (r[i].GetBegin() == r[j].GetBegin() && r[i].GetEnd() < r[j].GetEnd())
	})
}
