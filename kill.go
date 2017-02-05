package mesos

import (
	"time"
)

type KillOpt func(*KillPolicy)

func (f *KillPolicy) With(opts ...KillOpt) *KillPolicy {
	for _, o := range opts {
		o(f)
	}
	return f
}

func WithGracePeriod(d time.Duration) KillOpt {
	return func(f *KillPolicy) {
		f.GracePeriod = &DurationInfo{Nanoseconds: d.Nanoseconds()}
	}
}

func OptionalKillPolicy(fo ...KillOpt) *KillPolicy {
	if len(fo) == 0 {
		return nil
	}
	return (&KillPolicy{}).With(fo...)
}
