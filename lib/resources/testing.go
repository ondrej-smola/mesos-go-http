package resources

import "github.com/ondrej-smola/mesos-go-http/lib"

func Cpus(cpus float64) *mesos.Resource {
	return &mesos.Resource{
		Name: mesos.Strp(string(mesos.CPUS)),
		Type: mesos.Value_SCALAR.Enum(),
		Scalar: &mesos.Value_Scalar{
			Value: &cpus,
		},
	}
}

func Mem(mem float64) *mesos.Resource {
	return &mesos.Resource{
		Name: mesos.Strp(string(mesos.MEM)),
		Type: mesos.Value_SCALAR.Enum(),
		Scalar: &mesos.Value_Scalar{
			Value: &mem,
		},
	}
}

func Ports(ranges ...*mesos.Value_Range) *mesos.Resource {
	return &mesos.Resource{
		Name: mesos.Strp(string(mesos.PORTS)),
		Type: mesos.Value_RANGES.Enum(),
		Ranges: &mesos.Value_Ranges{
			Range: ranges,
		},
	}
}
