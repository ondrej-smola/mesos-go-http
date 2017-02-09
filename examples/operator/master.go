package main

import (
	"fmt"
	"github.com/ondrej-smola/mesos-go-http/client"
	"github.com/ondrej-smola/mesos-go-http/client/leader"
	"github.com/ondrej-smola/mesos-go-http/operator/master"
	"github.com/spf13/cobra"
	"os"
)

func masterCommands(cfg *config) *cobra.Command {
	getLeader := func(opts ...leader.Opt) *master.Client {
		return master.New(leader.New(cfg.endpoints))
	}

	var getHealth = &cobra.Command{
		Use: "health",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getLeader().GetHealth(cfg.ctx)
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
			v, err := getLeader().GetFlags(cfg.ctx)
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
			v, err := getLeader().GetVersion(cfg.ctx)
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
			v, err := getLeader().GetMetrics(&master.Call_GetMetrics{}, cfg.ctx)
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
			v, err := getLeader().GetLoggingLevel(cfg.ctx)
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
			v, err := getLeader().GetState(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getRoles = &cobra.Command{
		Use: "roles",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getLeader().GetRoles(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getWeights = &cobra.Command{
		Use: "weights",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getLeader().GetWeights(cfg.ctx)
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
			v, err := getLeader().GetFrameworks(cfg.ctx)
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
			v, err := getLeader().GetExecutors(cfg.ctx)
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
			v, err := getLeader().GetTasks(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getMaster = &cobra.Command{
		Use: "master",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getLeader().GetMaster(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var getQuota = &cobra.Command{
		Use: "quota",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getLeader().GetQuota(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var maintenance = &cobra.Command{
		Use: "maintenance",
	}

	var maintenanceStatus = &cobra.Command{
		Use: "status",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getLeader().GetMaintenanceStatus(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	var maintenanceSchedule = &cobra.Command{
		Use: "schedule",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getLeader().GetMaintenanceSchedule(cfg.ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.printResponse(v)
		},
	}

	maintenance.AddCommand(maintenanceStatus, maintenanceSchedule)

	var getEventStream = &cobra.Command{
		Use: "event-stream",
		Run: func(cmd *cobra.Command, args []string) {
			events := master.NewEventStream(
				leader.New(cfg.endpoints, leader.WithClientOpts(client.WithRecordIOFraming())),
				cfg.ctx,
			)
			for ev := range events {
				fmt.Println(ev)
				if ev.Err != nil {
					fmt.Println(ev.Err)
					os.Exit(1)
				} else {
					cfg.printResponse(ev.Event)
				}

			}
		},
	}

	var mCmd = &cobra.Command{
		Use: "master",
		PersistentPreRun: func(c *cobra.Command, args []string) {
			if len(cfg.endpoints) < 1 {
				panic("At least one endpoint is expected")
			}
		},
	}
	mCmd.AddCommand(getHealth, getFlags, getVersion, getMetrics)
	mCmd.AddCommand(getLoggingLevel, getState, getFrameworks)
	mCmd.AddCommand(getExecutors, getTasks, getRoles, getWeights)
	mCmd.AddCommand(getQuota, getMaster, maintenance)
	mCmd.AddCommand(getEventStream)
	return mCmd
}
