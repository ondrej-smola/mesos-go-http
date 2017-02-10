package callopt

import (
	"context"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
)

type (
	callOpts struct {
		opts []scheduler.CallOpt
		via  flow.Flow
	}
)

func Blueprint(opts ...scheduler.CallOpt) flow.StageBlueprint {
	return flow.StageBlueprintFunc(func(...flow.MatOpt) flow.Stage {
		return New(opts...)
	})
}

// Applies opts to all pushed scheduler calls
func New(opts ...scheduler.CallOpt) flow.Stage {
	return &callOpts{
		opts: opts,
	}
}

func (f *callOpts) Push(ev flow.Message, ctx context.Context) error {
	switch c := ev.(type) {
	case *scheduler.Call:
		for _, o := range f.opts {
			o(c)
		}
	}

	return f.via.Push(ev, ctx)
}

func (i *callOpts) Pull(ctx context.Context) (flow.Message, error) {
	return i.via.Pull(ctx)

}

func (i *callOpts) Via(f flow.Flow) {
	i.via = f
}

func (i *callOpts) Close() error {
	return i.via.Close()
}
