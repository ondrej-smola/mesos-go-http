package flow

import (
	"context"
	"github.com/go-kit/kit/log"
	"io"
)

type (
	// Mat config is used to modify flow just before materialization
	MatConfig struct {
		Log log.Logger
	}

	MatOpt func(c *MatConfig)

	MatOpts []MatOpt

	// Message is marker interface for representing entity that can be send through flow
	Message interface {
		IsFlowMessage()
	}

	// Flow chain of stage(s) and sink
	Flow interface {
		// Blocks until next message is ready or context is cancelled
		Pull(context.Context) (Message, error)
		// Blocks until message is processed or context is cancelled
		Push(Message, context.Context) error
		// Blocks until flow is closed.
		io.Closer
	}

	// Stage is connected to next stage or sink using Via func
	Stage interface {
		Flow
		Via(Flow)
	}

	// Represents terminal part of flow
	Sink interface {
		Flow
		IsSink()
	}

	// Blueprint can be materialized to create instance of flow
	Blueprint      func(...MatOpt) Flow
	StageBlueprint func(...MatOpt) Stage
	SinkBlueprint  func(...MatOpt) Sink

	// Flow builder builds up chain of stages that are terminated in sink
	FlowBuilder interface {
		Append(Stage) FlowBuilder
		RunWith(Sink) Flow
	}

	// Flow blueprint builder builds up chain of stages blueprints that are terminated in sink blueprint.
	// Connecting it with sink returns blueprint that can be used to create new flow instances
	FlowBlueprintBuilder interface {
		Append(StageBlueprint) FlowBlueprintBuilder
		RunWith(SinkBlueprint, ...MatOpt) Blueprint
	}

	// Implementation of FlowBuilder
	flowBuilder struct {
		processors []Stage
		sink       Sink
	}

	// Implementation of FlowBlueprintBuilder
	blueprintBuilder struct {
		processors []StageBlueprint
		sink       SinkBlueprint
	}
)

// Materialize flow based on blueprint
func (f Blueprint) Mat(opts ...MatOpt) Flow {
	return f(opts...)
}

// Materialize stage based on blueprint
func (f StageBlueprint) Mat(opts ...MatOpt) Stage {
	return f(opts...)
}

// Materialize stage based on blueprint
func (f SinkBlueprint) Mat(opts ...MatOpt) Sink {
	return f(opts...)
}

func New(p ...Stage) FlowBuilder {
	return &flowBuilder{processors: p}
}

func (l *flowBuilder) Append(p Stage) FlowBuilder {
	l.processors = append(l.processors, p)
	return l
}

func (l *flowBuilder) RunWith(s Sink) Flow {
	if len(l.processors) > 0 {
		for i := 1; i < len(l.processors); i++ {
			l.processors[i-1].Via(l.processors[i])
		}

		l.processors[len(l.processors)-1].Via(s)
		return l.processors[0]
	} else {
		return s
	}
}

func BlueprintBuilder(p ...StageBlueprint) FlowBlueprintBuilder {
	return &blueprintBuilder{processors: p}
}

func (l *blueprintBuilder) Append(p StageBlueprint) FlowBlueprintBuilder {
	l.processors = append(l.processors, p)
	return l
}

func (l *blueprintBuilder) RunWith(s SinkBlueprint, opts ...MatOpt) Blueprint {
	return func(addOpts ...MatOpt) Flow {
		allOpts := append(opts, addOpts...)

		f := New()
		for _, pb := range l.processors {
			f.Append(pb.Mat(allOpts...))
		}
		return f.RunWith(s.Mat(allOpts...))
	}
}

// Set logger for flow, flow components should inherit this logger
func WithLogger(l log.Logger) MatOpt {
	return func(c *MatConfig) {
		c.Log = l
	}
}

// Utility function for converting slice of MatOpt to MatConfig
func (opts MatOpts) Config() *MatConfig {
	cfg := &MatConfig{}
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}
