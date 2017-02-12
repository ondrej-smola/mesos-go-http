package main

import "github.com/spf13/cobra"

type config struct {
	endpoints []string

	taskCpus  float64
	taskMem   float64
	taskImage string
	taskCmd   string
	taskArgs  []string
	numTasks  int
}

func main() {

	cfg := &config{}

	var rootCmd = &cobra.Command{
		Use: "Mesos Operator CLI",
	}

	rootCmd.Flags().StringSliceVarP(
		&cfg.endpoints, "endpoints", "e", []string{"http://127.0.0.1:5050/api/v1/scheduler"}, "operator api endpoint",
	)

	rootCmd.Flags().Float64Var(&cfg.taskCpus, "cpus", 0.1, "task cpus")
	rootCmd.Flags().Float64Var(&cfg.taskMem, "mem", 0.1, "task mem")
	rootCmd.Flags().StringVar(&cfg.taskImage, "image", "alpine:3.4", "task image")
	rootCmd.Flags().StringVar(&cfg.taskCmd, "cmd", "echo", "task cmd")
	rootCmd.Flags().StringSliceVar(&cfg.taskArgs, "arg", []string{"hello"}, "task cmd argument")
	rootCmd.Flags().IntVar(&cfg.numTasks, "tasks", 5, "total number of tasks to launch")

	rootCmd.Run = func(c *cobra.Command, args []string) {
		run(cfg)
	}

	rootCmd.Execute()
}
