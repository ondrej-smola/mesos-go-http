package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/client"
	"github.com/ondrej-smola/mesos-go-http/client/leader"
	"github.com/ondrej-smola/mesos-go-http/codec"
	"github.com/ondrej-smola/mesos-go-http/operator/master"
	"github.com/pkg/errors"
	"os"
)

func main() {

	w := log.NewSyncWriter(os.Stderr)
	logger := log.NewContext(log.NewLogfmtLogger(w))

	c := leader.New(
		leader.WithMasters(os.Args[1:]...),
		leader.WithLogger(logger),
		leader.WithClientProvider(client.NewProvider(client.WithCodec(codec.JsonCodec))),
	)

	resp, err := c.Do(&master.Call{
		Type: master.Call_GET_AGENTS,
	}, context.Background())

	if err != nil {
		panic(errors.Wrap(err, "Client do"))
	}
	defer resp.Close()

	ev := &master.Response{}
	err = resp.Read(ev)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	if ev.Type != master.Response_GET_AGENTS {
		panic("Unexpected response type: " + ev.Type.String())
	}

	for _, a := range ev.GetAgents.Agents {
		fmt.Println(fmt.Sprintf("Agent %v, resources %v", a.AgentInfo.ID.Value, mesos.Resources(a.TotalResources)))
	}
}
