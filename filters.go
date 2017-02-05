package mesos

import (
	"math/rand"
	"time"
)

type FilterOpt func(*Filters)

func (f *Filters) With(opts ...FilterOpt) *Filters {
	for _, o := range opts {
		o(f)
	}
	return f
}

func RefuseSeconds(d time.Duration) FilterOpt {
	return func(f *Filters) {
		s := d.Seconds()
		f.RefuseSeconds = &s
	}
}

func RefuseSecondsWithJitter(d time.Duration) FilterOpt {
	return func(f *Filters) {
		s := time.Duration(rand.Int63n(int64(d))).Seconds()
		f.RefuseSeconds = &s
	}
}

func OptionalFilters(fo ...FilterOpt) *Filters {
	if len(fo) == 0 {
		return nil
	}
	return (&Filters{}).With(fo...)
}
