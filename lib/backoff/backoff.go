package backoff

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type (
	Opt func(c *Config)

	Provider interface {
		New(ctx ...context.Context) Backoff
	}

	Backoff interface {
		Attempts() <-chan int
		Close()
		Reset()
	}

	expBackoff struct {
		cfg    *Config
		c      <-chan int
		reset  chan struct{}
		ctx    context.Context
		cancel context.CancelFunc
	}

	Config struct {
		rnd               RandFunc
		maxAttempts       int
		minWait           time.Duration
		maxWait           time.Duration
		backoffFactor     float64
		jitterMaxFraction float64 // add to duration up to a fraction of current duration
	}

	// must return number between [0,1)
	RandFunc func() float64
)

func Always() Opt {
	return WithMaxAttempts(math.MaxInt64)
}

func Once() Opt {
	return WithMaxAttempts(1)
}

func WithMaxAttempts(attempts int) Opt {
	if attempts <= 0 {
		panic(fmt.Sprintf("Attempts must be > 0, got %v", attempts))
	}

	return func(c *Config) {
		c.maxAttempts = attempts
	}
}

func WithMinWait(w time.Duration) Opt {
	if w <= 0 {
		panic(fmt.Sprintf("retry minWait must be positive got %v", w))
	}

	return func(c *Config) {
		c.minWait = w
	}
}

func WithMaxWait(w time.Duration) Opt {
	if w <= 0 {
		panic(fmt.Sprintf("retry maxWait must be positive got %v", w))
	}

	return func(c *Config) {
		c.maxWait = w
	}
}

func WithBackoffFactor(f float64) Opt {
	if f < 1.0 {
		panic(fmt.Sprintf("retry backoff factor must be >= 1 - got %v", f))
	}

	return func(c *Config) {
		c.backoffFactor = f
	}
}

func WithRandFunc(f RandFunc) Opt {
	return func(c *Config) {
		c.rnd = f
	}
}

// Up to what fraction of current backoff interval can be added as jitter
func WithJitterFraction(f float64) Opt {
	if f < 0 || f > 1 {
		panic(fmt.Sprintf("jitter fraction must be between [0,1] - got %v", f))
	}

	return func(c *Config) {
		c.jitterMaxFraction = f
	}
}

func New(opts ...Opt) Provider {
	cfg := &Config{
		rnd:               rand.Float64,
		maxAttempts:       5,
		minWait:           1 * time.Second,
		maxWait:           15 * time.Second,
		backoffFactor:     2,
		jitterMaxFraction: 0.2,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.maxWait < cfg.minWait {
		panic(fmt.Sprintf("retry maxWait (%v) < cfg.minWait (%v)", cfg.maxWait, cfg.minWait))
	}

	return cfg
}

// creates new backoff with optional Context
func (c *Config) New(optionalCtx ...context.Context) Backoff {
	ctx := context.Background()

	if len(optionalCtx) > 1 {
		panic("Only up to one optional context must be provided")
	}

	if len(optionalCtx) > 0 {
		ctx = optionalCtx[0]
	}

	ctx, cancel := context.WithCancel(ctx)
	reset := make(chan struct{})
	e := &expBackoff{
		cfg:    c,
		cancel: cancel,
		ctx:    ctx,
		reset:  reset,
	}

	e.c = e.run()

	return e
}

func (a *expBackoff) Attempts() <-chan int {
	return a.c
}

func (a *expBackoff) Close() {
	a.cancel()
}

func (a *expBackoff) Reset() {
	select {
	case a.reset <- struct{}{}:
	case <-a.ctx.Done():
	}
}

func (e *expBackoff) run() <-chan int {
	cfg := e.cfg
	tokens := make(chan int)
	limiter := tokens
	go func() {
		defer close(tokens)
		t := time.NewTimer(cfg.minWait)
		defer t.Stop()
		attempt := 1
		for {
			select {
			case limiter <- attempt:
				if attempt == cfg.maxAttempts {
					return
				} else {
					mul := math.Pow(float64(cfg.backoffFactor), float64(attempt-1))
					next := time.Duration(float64(cfg.minWait) * mul)

					jitter := cfg.rnd() * cfg.jitterMaxFraction * float64(next)
					next += time.Duration(jitter)

					if next < cfg.minWait {
						next = cfg.minWait
					} else if next > cfg.maxWait {
						next = cfg.maxWait
					}
					limiter = nil

					t.Reset(next)
					attempt += 1
				}
			case <-e.reset:
				attempt = 1
				t.Reset(0)
			case <-t.C:
				limiter = tokens
			case <-e.ctx.Done():
				return
			}
		}
	}()
	return tokens
}
