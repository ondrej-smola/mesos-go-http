package main

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http/client"
	"github.com/ondrej-smola/mesos-go-http/operator/agent"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	endpoint := ""
	marsh := jsonpb.Marshaler{}

	getAgent := func() *agent.Client {
		return agent.New(client.New(endpoint))
	}

	printResponse := func(c proto.Message) {
		b, err := marsh.MarshalToString(c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println(b)
		}
	}

	var getHealth = &cobra.Command{
		Use: "health",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetHealth(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printResponse(v)
		},
	}

	var getFlags = &cobra.Command{
		Use: "flags",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetFlags(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printResponse(v)
		},
	}

	var getVersion = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetVersion(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printResponse(v)
		},
	}

	var getMetrics = &cobra.Command{
		Use: "metrics",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetMetrics(&agent.Call_GetMetrics{}, context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printResponse(v)
		},
	}

	var getLoggingLevel = &cobra.Command{
		Use: "logging-level",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetLoggingLevel(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printResponse(v)
		},
	}

	var getState = &cobra.Command{
		Use: "state",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetState(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printResponse(v)
		},
	}

	var getContainers = &cobra.Command{
		Use: "containers",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetContainers(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			printResponse(v)
		},
	}

	var getFrameworks = &cobra.Command{
		Use: "frameworks",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetFrameworks(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printResponse(v)
		},
	}

	var getExecutors = &cobra.Command{
		Use: "executors",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetExecutors(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printResponse(v)
		},
	}

	var getTasks = &cobra.Command{
		Use: "tasks",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := getAgent().GetTasks(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printResponse(v)
		},
	}

	var rootCmd = &cobra.Command{Use: "Mesos Operator Example Agent CLI"}
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "http://127.0.0.1:5051/api/v1", "host port of agent")
	rootCmd.AddCommand(getHealth, getFlags, getVersion, getMetrics)
	rootCmd.AddCommand(getLoggingLevel, getState, getContainers, getFrameworks)
	rootCmd.AddCommand(getExecutors, getTasks)
	rootCmd.Execute()
}
