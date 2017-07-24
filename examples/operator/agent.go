package main

import (
	"fmt"
	"os"

	"github.com/ondrej-smola/mesos-go-http/lib/client"
	"github.com/ondrej-smola/mesos-go-http/lib/operator/agent"
	"github.com/spf13/cobra"
)

func agentCommands(cfg *config) *cobra.Command {

	getAgent := func() *agent.Client {
		return agent.New(client.New(cfg.endpoints[0]))
	}

	var getHealth = &cobra.Command{
		Use: "health",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetHealth(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getFlags = &cobra.Command{
		Use: "flags",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetFlags(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getVersion = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetVersion(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getMetrics = &cobra.Command{
		Use: "metrics",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetMetrics(&agent.Call_GetMetrics{}, cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getLoggingLevel = &cobra.Command{
		Use: "logging-level",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetLoggingLevel(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getState = &cobra.Command{
		Use: "state",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetState(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getContainers = &cobra.Command{
		Use: "containers",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetContainers(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			cfg.printResponse(v)
		},
	}

	var getFrameworks = &cobra.Command{
		Use: "frameworks",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetFrameworks(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getExecutors = &cobra.Command{
		Use: "executors",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetExecutors(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getTasks = &cobra.Command{
		Use: "tasks",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetTasks(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var agentCmd = &cobra.Command{
		Use: "agent",
		PersistentPreRun: func(c *cobra.Command, args []string) {
			if len(cfg.endpoints) != 1 {
				panic("Exactly one endpoint is expected")
			}
		},
	}
	agentCmd.AddCommand(getHealth, getFlags, getVersion, getMetrics)
	agentCmd.AddCommand(getLoggingLevel, getState, getContainers, getFrameworks)
	agentCmd.AddCommand(getExecutors, getTasks)
	return agentCmd
}
