package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type PrometheusMetrics struct {
	pushPull        *prometheus.CounterVec
	pushPullErr     *prometheus.CounterVec
	pushPullLatency *prometheus.HistogramVec
	offers          *prometheus.CounterVec
	resources       *prometheus.CounterVec
}

func New() *PrometheusMetrics {
	p := &PrometheusMetrics{}

	p.pushPull = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "examples",
		Subsystem: "scheduler",
		Name:      "push_pull",
		Help:      "msg push/pull count",
	}, []string{"type", "name"})

	p.pushPullErr = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "examples",
		Subsystem: "scheduler",
		Name:      "push_pull_err",
		Help:      "msg push/pull errors count",
	}, []string{"type"})

	p.pushPullLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "examples",
		Subsystem: "scheduler",
		Name:      "push_pull_latency",
		Help:      "msg push/pull latency",
	}, []string{"type"})

	p.offers = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "examples",
		Subsystem: "scheduler",
		Name:      "offers",
		Help:      "recived/declined offers count",
	}, []string{"type"})

	p.resources = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "examples",
		Subsystem: "scheduler",
		Name:      "resources",
		Help:      "offered resources count",
	}, []string{"type"})

	return p
}

func (t *PrometheusMetrics) Push(name string) {
	t.pushPull.WithLabelValues("push", name).Add(1)
}

func (t *PrometheusMetrics) Pull(name string) {
	t.pushPull.WithLabelValues("pull", name).Add(1)
}

func (t *PrometheusMetrics) PushErr(err error) {
	t.pushPullErr.WithLabelValues("push").Add(1)
}

func (t *PrometheusMetrics) PullErr(err error) {
	t.pushPullErr.WithLabelValues("pull", "name").Add(1)
}

func (t *PrometheusMetrics) PushLatency(time time.Duration) {
	t.pushPullLatency.WithLabelValues("push").Observe(time.Seconds())
}

func (t *PrometheusMetrics) PullLatency(time time.Duration) {
	t.pushPullLatency.WithLabelValues("pull").Observe(time.Seconds())
}

func (t *PrometheusMetrics) OffersReceived(count uint32) {
	t.offers.WithLabelValues("received").Add(float64(count))
}

func (t *PrometheusMetrics) OffersDeclined(count uint32) {
	t.offers.WithLabelValues("declined").Add(float64(count))
}

func (t *PrometheusMetrics) ResourceOffered(name string, value float64) {
	t.resources.WithLabelValues(name).Add(float64(value))
}

func (p *PrometheusMetrics) MustRegister(r *prometheus.Registry) {
	r.MustRegister(p.pushPull)
	r.MustRegister(p.pushPullErr)
	r.MustRegister(p.pushPullLatency)
	r.MustRegister(p.offers)
	r.MustRegister(p.resources)
}
