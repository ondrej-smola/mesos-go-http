package mesos

import "github.com/gogo/protobuf/proto"

func (f *Filters) Clone() *Filters {
	return proto.Clone(f).(*Filters)
}

func (p *KillPolicy) Clone() *KillPolicy {
	return proto.Clone(p).(*KillPolicy)
}

func (r *Resource) Clone() *Resource {
	return proto.Clone(r).(*Resource)
}

func (r *Value_Range) Clone() *Value_Range {
	return &Value_Range{Begin: UI64p(r.GetBegin()), End: UI64p(r.GetEnd())}
}

func (resources Resources) Clone() Resources {
	res := make(Resources, len(resources))
	for i, r := range resources {
		res[i] = r.Clone()
	}

	return res
}

func (r Ranges) Clone() Ranges {
	ret := make(Ranges, len(r))

	for i := range r {
		ret[i] = r[i].Clone()
	}

	return ret
}
