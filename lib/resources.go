package mesos

import (
	"bytes"
	"strconv"
)

type Resources []*Resource

type ResourceName string

const (
	CPUS  = ResourceName("cpus")
	MEM   = ResourceName("mem")
	PORTS = ResourceName("ports")
	GPUS  = ResourceName("gpus")
	DISK  = ResourceName("disk")
)

func (r *Resource) ScalarValueOrZero() float64 {
	if s := r.GetScalar(); s != nil {
		return s.GetValue()
	} else {
		return 0
	}
}

func (r *Resource) RangesOrZero() Ranges {
	if r := r.GetRanges(); r != nil {
		return r.Range
	} else {
		return nil
	}
}

func (r *Resource) WithRole(role string) *Resource {
	r.Role = &role
	return r
}

func (resources Resources) String() string {
	if len(resources) == 0 {
		return ""
	}
	buf := bytes.Buffer{}
	for i, r := range resources {
		if i > 0 {
			buf.WriteString(";")
		}

		buf.WriteString(r.GetName())
		buf.WriteString("(")
		buf.WriteString(r.GetRole())
		if ri := r.GetReservation(); ri != nil {
			buf.WriteString(", ")
			buf.WriteString(ri.GetPrincipal())
		}
		buf.WriteString(")")
		if d := r.GetDisk(); d != nil {
			buf.WriteString("[")
			if p := d.GetPersistence(); p != nil {
				buf.WriteString(p.GetId())
			}
			if v := d.GetVolume(); v != nil {
				buf.WriteString(":")
				vconfig := v.GetContainerPath()
				if h := v.GetHostPath(); h != "" {
					vconfig = h + ":" + vconfig
				}
				if m := v.Mode; m != nil {
					switch *m {
					case Volume_RO:
						vconfig += ":ro"
					case Volume_RW:
						vconfig += ":rw"
					}
				}
				buf.WriteString(vconfig)
			}
			buf.WriteString("]")
		}
		buf.WriteString(":")

		switch r.GetType() {
		case Value_SCALAR:
			buf.WriteString(strconv.FormatFloat(r.GetScalar().GetValue(), 'f', -1, 64))
		case Value_RANGES:
			buf.WriteString("[")
			for j, r := range r.RangesOrZero() {
				if j > 0 {
					buf.WriteString(",")
				}
				buf.WriteString(strconv.FormatUint(r.GetBegin(), 10))
				buf.WriteString("-")
				buf.WriteString(strconv.FormatUint(r.GetEnd(), 10))
			}
			buf.WriteString("]")
		case Value_SET:
			buf.WriteString("{")
			items := r.GetSet().GetItem()
			for j := range items {
				if j > 0 {
					buf.WriteString(",")
				}
				buf.WriteString(items[j])
			}
			buf.WriteString("}")
		}
	}
	return buf.String()
}
