package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/ondrej-smola/mesos-go-http/lib"
	"github.com/ondrej-smola/mesos-go-http/lib/backoff"
	"github.com/ondrej-smola/mesos-go-http/lib/client"
	"github.com/ondrej-smola/mesos-go-http/lib/client/leader"
	"github.com/ondrej-smola/mesos-go-http/examples/scheduler/metrics"
	"github.com/ondrej-smola/mesos-go-http/lib/flow"
	"github.com/ondrej-smola/mesos-go-http/lib/resources/filter"
	"github.com/ondrej-smola/mesos-go-http/lib/resources/find"
	"github.com/ondrej-smola/mesos-go-http/lib/scheduler"
	"github.com/ondrej-smola/mesos-go-http/lib/scheduler/stage/ack"
	"github.com/ondrej-smola/mesos-go-http/lib/scheduler/stage/fwid"
	"github.com/ondrej-smola/mesos-go-http/lib/scheduler/stage/heartbeat"
	"github.com/ondrej-smola/mesos-go-http/lib/scheduler/stage/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type app struct {
	cfg           *config
	tasksLaunched int
	log           log.Logger
}

func run(cfg *config) {
	// Setup
	w := log.NewSyncWriter(os.Stderr)
	logger := log.NewContext(log.NewLogfmtLogger(w)).With("ts", log.DefaultTimestampUTC)

	metricsBackend := metrics.New()

	sched := scheduler.Blueprint(
		leader.New(
			cfg.endpoints,
			leader.WithLogger(log.NewContext(logger).With("src", "leader_client")),
			leader.WithClientOpts(client.WithRecordIOFraming()),
		),
	)

	blueprint := flow.BlueprintBuilder().
		Append(monitor.Blueprint(metricsBackend)).
		Append(heartbeat.Blueprint()).
		Append(ack.Blueprint()).
		Append(fwid.Blueprint()).
		RunWith(sched, flow.WithLogger(log.NewContext(logger).With("src", "flow")))
	//

	go serveMetrics(metricsBackend, cfg.metricsBind, logger)

	a := &app{
		cfg: cfg,
		log: log.NewContext(logger).With("src", "main"),
	}

	ctx := context.Background()

	retry := backoff.New(backoff.Always()).New(ctx)
	defer retry.Close()

	var msg flow.Message

	subscribeCall := scheduler.Subscribe(&mesos.FrameworkInfo{
		User: mesos.Strp("root"),
		Name: mesos.Strp("test"),
	})
	
	for attempt := range retry.Attempts() {
		a.tasksLaunched = 0
		a.log.Log("event", "connecting", "attempt", attempt)
		fl := blueprint.Mat()
		err := fl.Push(subscribeCall, ctx)
		for err == nil {
			msg, err = fl.Pull(ctx)
			if err == nil {
				switch m := msg.(type) {
				case *scheduler.Event:
					a.log.Log("event", "message_received", "type", m.Type.String())
					switch m.GetType() {
					case scheduler.Event_SUBSCRIBED:
						retry.Reset()
					case scheduler.Event_UPDATE:
						status := m.Update.Status
						a.log.Log(
							"event", "status_update",
							"task_id", status.TaskId.GetValue(),
							"status", status.State.String(),
							"msg", status.Message,
						)
					case scheduler.Event_OFFERS:
						err = a.handleOffers(m.Offers.Offers, fl)
					}
				}
			}
		}

		a.log.Log("event", "failed", "attempt", attempt, "err", err)
		fl.Close()
	}
}
func serveMetrics(metrics *metrics.PrometheusMetrics, endpoint string, log log.Logger) {
	if endpoint == "" {
		log.Log("event", "metrics_server_disable")
		return
	}

	promRegistry := prometheus.NewRegistry()
	metrics.MustRegister(promRegistry)

	l, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Log("event", "http_listener", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))

	srv := http.Server{
		Addr:    endpoint,
		Handler: mux,
	}

	log.Log("event", "metrics_server", "listening", l.Addr().String())
	if err := srv.Serve(l); err != nil {
		log.Log("event", "metrics_server", "err", err)
		os.Exit(1)
	}
}

func (a *app) handleOffers(offers []*mesos.Offer, flow flow.Flow) error {
	useShell := false

	for _, o := range offers {
		logger := log.NewContext(a.log).With("offer", o.Id.GetValue())

		resources := mesos.Resources(o.Resources)
		logger.Log("resources", resources)
		tasks := []*mesos.TaskInfo{}

		availableResources := filter.All(filter.Unreserved(), resources...)

		for a.tasksLaunched < a.cfg.numTasks {
			logger.Log("resources", availableResources)

			cpus, remaining, ok := find.Scalar(mesos.CPUS, a.cfg.taskCpus, availableResources...)
			if !ok {
				break
			}
			mem, remaining, ok := find.Scalar(mesos.MEM, a.cfg.taskMem, remaining...)
			if !ok {
				break
			}

			// remaining resources are available for next task
			availableResources = remaining

			a.tasksLaunched++
			taskID := a.tasksLaunched

			t := &mesos.TaskInfo{
				Name:    mesos.Strp("Task " + strconv.Itoa(taskID)),
				TaskId:  &mesos.TaskID{Value: mesos.Strp(strconv.Itoa(taskID))},
				SlaveId: o.SlaveId,
				Container: &mesos.ContainerInfo{
					Type: mesos.ContainerInfo_DOCKER.Enum(),
					Docker: &mesos.ContainerInfo_DockerInfo{
						Image:   mesos.Strp(a.cfg.taskImage),
						Network: mesos.ContainerInfo_DockerInfo_NONE.Enum(),
					},
				},
				Command: &mesos.CommandInfo{
					Shell:     &useShell,
					Value:     nilIfEmptyString(a.cfg.taskCmd),
					Arguments: a.cfg.taskArgs,
				},
				Resources: []*mesos.Resource{cpus, mem},
			}

			tasks = append(tasks, t)
		}

		if len(tasks) > 0 {
			logger.Log("event", "launch", "count", len(tasks))

			err := flow.Push(scheduler.AcceptOffer(o.Id,scheduler.OpLaunch(tasks...)), context.Background())
			if err != nil {
				return err
			}
		} else {
			logger.Log("event", "declined")
			err := flow.Push(scheduler.Decline(o.Id), context.Background())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func nilIfEmptyString(cmd string) *string {
	if cmd == "" {
		return nil
	} else {
		return &cmd
	}
}
