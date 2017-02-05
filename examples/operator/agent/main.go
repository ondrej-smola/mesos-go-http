package main

import (
	"context"
	"fmt"
	"github.com/ondrej-smola/mesos-go-http/client"
	"github.com/ondrej-smola/mesos-go-http/operator/agent"
	"github.com/pkg/errors"
)

func main() {
	c := client.NewProvider().New(agent.EndpointFunc("10.0.75.2:5061"))

	resp, err := c.Do(&agent.Call{
		Type: agent.Call_GET_VERSION,
	}, context.Background())

	if err != nil {
		panic(errors.Wrap(err, "Client do"))
	}
	defer resp.Close()

	ev := &agent.Response{}
	err = resp.Read(ev)
	if err != nil {
		panic(errors.Wrap(err, "Response read"))
	}

	fmt.Println(ev.GetVersion)
}
