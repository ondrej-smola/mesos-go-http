package resources

import "github.com/ondrej-smola/mesos-go-http/lib"

func Cpus(cpus float64) *mesos.Resource {
	return NewScalar(mesos.CPUS, cpus)
}

func Mem(mem float64) *mesos.Resource {
	return NewScalar(mesos.MEM, mem)
}

func Ports(ranges ...*mesos.Value_Range) *mesos.Resource {
	return NewRange(mesos.PORTS, ranges...)
}

func NewScalar(name mesos.ResourceName, value float64) *mesos.Resource {
	return &mesos.Resource{
		Name: mesos.Strp(string(name)),
		Type: mesos.Value_SCALAR.Enum(),
		Scalar: &mesos.Value_Scalar{
			Value: &value,
		},
	}
}

func NewRange(name mesos.ResourceName, ranges ...*mesos.Value_Range) *mesos.Resource {
	return &mesos.Resource{
		Name: mesos.Strp(string(name)),
		Type: mesos.Value_RANGES.Enum(),
		Ranges: &mesos.Value_Ranges{
			Range: ranges,
		},
	}
}
