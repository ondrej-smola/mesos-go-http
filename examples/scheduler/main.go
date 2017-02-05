package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/backoff"
	"github.com/ondrej-smola/mesos-go-http/client/leader"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
	"github.com/ondrej-smola/mesos-go-http/scheduler/stage/ack"
	"github.com/ondrej-smola/mesos-go-http/scheduler/stage/callopt"
	"github.com/ondrej-smola/mesos-go-http/scheduler/stage/fwid"
	"github.com/ondrej-smola/mesos-go-http/scheduler/stage/heartbeat"
	"os"
	"strconv"
	"time"
)

type app struct {
	totalTasks    int
	tasksLaunched int
	wants         mesos.Resources
	log           log.Logger
}

func main() {
	// Setup
	w := log.NewSyncWriter(os.Stderr)
	logger := log.NewContext(log.NewLogfmtLogger(w)).With("ts", log.DefaultTimestampUTC)

	masters := mesos.NewMasters("10.0.75.2:5050", "10.0.75.2:5051", "10.0.75.2:5052")

	sched := scheduler.Blueprint(
		leader.New(
			scheduler.EndpointFunc,
			masters,
			leader.WithLogger(log.NewContext(logger).With("component", "leader_client"))),
	)

	blueprint := flow.BlueprintBuilder().
		Append(callopt.Blueprint(scheduler.Filters(mesos.RefuseSecondsWithJitter(3*time.Second)))).
		Append(heartbeat.Blueprint()).
		Append(ack.Blueprint()).
		Append(fwid.Blueprint()).
		RunWith(sched, flow.WithLogger(log.NewContext(logger).With("component", "flow")))
	//

	wants := mesos.Resources{}
	wants = wants.Add(
		mesos.BuildResource().Name("cpus").Scalar(1).Build(),
		mesos.BuildResource().Name("mem").Scalar(256).Build(),
	)
	a := &app{
		totalTasks: 5,
		wants:      wants,
		log:        log.NewContext(logger).With("component", "main"),
	}

	ctx := context.Background()

	retry := backoff.New(backoff.Always()).New(ctx)
	defer retry.Close()

	var msg flow.Message

	for attempt := range retry.Attempts() {
		a.tasksLaunched = 0
		a.log.Log("action", "connecting", "attempt", attempt)
		fl := blueprint.Mat()
		err := fl.Push(scheduler.Subscribe(mesos.FrameworkInfo{User: "root", Name: "test"}), ctx)
		for err == nil {
			msg, err = fl.Pull(ctx)
			if err == nil {
				switch m := msg.(type) {
				case *scheduler.Event:
					a.log.Log("action", "message_received", "type", m.Type.String())
					switch m.Type {
					case scheduler.Event_SUBSCRIBED:
						retry.Reset()
					case scheduler.Event_UPDATE:
						status := m.Update.Status
						a.log.Log(
							"action", "state_update",
							"task", status.TaskID.Value,
							"status", status.State.String(),
							"msg", status.Message,
						)
					case scheduler.Event_OFFERS:
						err = a.handleOffers(m.Offers.Offers, fl)
					}
				}
			}
		}

		a.log.Log("state", "failed", "attempt", attempt, "err", err)
	}
}

func (a *app) handleOffers(offers []mesos.Offer, flow flow.Flow) error {
	useShell := false

	for _, o := range offers {
		logger := log.NewContext(a.log).With("offer", o.ID.Value)

		offerResources := mesos.Resources(o.Resources)
		logger.Log("resources", offerResources)
		tasks := []mesos.TaskInfo{}

		for a.tasksLaunched < a.totalTasks && offerResources.ContainsAll(a.wants) {
			a.tasksLaunched++
			taskID := a.tasksLaunched

			logger.Log("action", "append_task", "task_id", taskID)

			t := mesos.TaskInfo{
				Name:    "Task " + strconv.Itoa(taskID),
				TaskID:  mesos.TaskID{Value: strconv.Itoa(taskID)},
				AgentID: o.AgentID,
				Container: &mesos.ContainerInfo{
					Type: mesos.ContainerInfo_DOCKER.Enum(),
					Docker: &mesos.ContainerInfo_DockerInfo{
						Image:   "nginx",
						Network: mesos.ContainerInfo_DockerInfo_NONE.Enum(),
					},
				},
				Command: &mesos.CommandInfo{
					Shell: &useShell,
				},
				Resources: offerResources.Find(a.wants),
			}

			tasks = append(tasks, t)
			offerResources = offerResources.Subtract(t.Resources...)
		}

		if len(tasks) > 0 {
			logger.Log("action", "launch_tasks", "count", len(tasks))
			accept := scheduler.Accept(
				scheduler.OfferOperations{scheduler.OpLaunch(tasks...)}.WithOffers(o.ID),
			)
			return flow.Push(accept, context.Background())
		} else {
			logger.Log("action", "decline_offer")
			return flow.Push(scheduler.Decline(o.ID), context.Background())
		}
	}

	return nil
}
