package filter

import (
	"github.com/ondrej-smola/mesos-go-http/lib"
)

type (
	Filter interface {
		Accepts(*mesos.Resource) bool
	}

	FilterFunc func(*mesos.Resource) bool
)

var _ = Filter(FilterFunc(nil))

func (f FilterFunc) Accepts(r *mesos.Resource) bool {
	return f(r)
}

func Role(role string) Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		return r.GetRole() == role
	})
}

func Scalar() Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		return r.GetType() == mesos.Value_SCALAR
	})
}

func Range() Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		return r.GetType() == mesos.Value_RANGES
	})
}

func Set() Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		return r.GetType() == mesos.Value_SET
	})
}

func Text() Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		return r.GetType() == mesos.Value_TEXT
	})
}

func Name(name mesos.ResourceName) Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		return r.GetName() == string(name)
	})
}

func Cpus() Filter {
	return Name(mesos.CPUS)
}

func WithCpus(cpus float64) Filter {
	return MinScalarValue(mesos.CPUS, cpus)
}

func MinScalarValue(name mesos.ResourceName, value float64) Filter {
	return And(Name(name), Scalar(), FilterFunc(func(r *mesos.Resource) bool {
		return r.ScalarValueOrZero() >= value
	}))
}

func Mem() Filter {
	return Name(mesos.MEM)
}

func WithMem(mem float64) Filter {
	return MinScalarValue(mesos.MEM, mem)
}

func Ports() Filter {
	return Name(mesos.PORTS)
}

func Gpus() Filter {
	return Name(mesos.GPUS)
}

func WithGpus(gpus float64) Filter {
	return MinScalarValue(mesos.GPUS, gpus)
}

func Disk() Filter {
	return Name(mesos.DISK)
}

func Revocable() Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		return r.Revocable != nil
	})
}

func Unreserved() Filter {
	return Role(mesos.Default_Resource_Role)
}

func StaticallyReserved(r *mesos.Resource) Filter {
	and := And(Not(Unreserved()), Not(DynamicallyReserved()))

	return FilterFunc(func(r *mesos.Resource) bool {
		return and.Accepts(r)
	})
}

func DynamicallyReserved() Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		return r.GetReservation() != nil
	})
}

func Or(filters ...Filter) Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		for _, f := range filters {
			if f.Accepts(r) {
				return true
			}
		}

		return false
	})
}

func Not(f Filter) Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		return !f.Accepts(r)
	})
}

func And(filters ...Filter) Filter {
	return FilterFunc(func(r *mesos.Resource) bool {
		for _, f := range filters {
			if !f.Accepts(r) {
				return false
			}
		}

		return true
	})
}

func All(rf Filter, resources ...*mesos.Resource) (res mesos.Resources) {
	for i := range resources {
		if rf.Accepts(resources[i]) {
			res = append(res, resources[i])
		}
	}

	return mesos.Resources(res)
}

func First(rf Filter, resources ...*mesos.Resource) (found *mesos.Resource, rest mesos.Resources) {
	for i, r := range resources {
		if found == nil && rf.Accepts(r) {
			found = resources[i]
		} else {
			rest = append(rest, resources[i])
		}
	}

	return
}
