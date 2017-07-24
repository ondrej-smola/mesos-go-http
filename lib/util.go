package mesos

import (
	"math"
)

func SubtractFixed(a, b float64) float64 {
	diff := ToFixed64(a) - ToFixed64(b)
	return FixedToFloat64(diff)
}

func AddFixed(a, b float64) float64 {
	return FixedToFloat64(ToFixed64(a) + ToFixed64(b))
}

func FixedToFloat64(f int64) float64 {
	// NOTE: We do the conversion from fixed point via integer division
	// and then modulus, rather than a single floating point division.
	// This ensures that we only apply floating point division to inputs
	// in the range [0,999], which is easier to check for correctness.
	var (
		quotient  = float64(f / 1000)
		remainder = float64(f%1000) / 1000.0
	)
	return quotient + remainder
}

func ToFixed64(f float64) int64 {
	return round64(f * 1000)
}

func round64(f float64) int64 {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int64(f + math.Copysign(0.5, f))
}

func F64p(f float64) *float64 {
	return &f
}

func F32p(f float32) *float32 {
	return &f
}

func UI64p(i uint64) *uint64 {
	return &i
}

func UI32p(i uint32) *uint32 {
	return &i
}

func I64p(i int64) *int64 {
	return &i
}

func I32p(i int32) *int32 {
	return &i
}

func BoolP(b bool) *bool {
	return &b
}

func Strp(s string) *string {
	return &s
}
