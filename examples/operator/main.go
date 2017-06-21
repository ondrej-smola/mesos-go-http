package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
)

type config struct {
	endpoints     []string
	printResponse func(proto.Message)
	ctx           context.Context
}

func main() {
	cfg := &config{
		printResponse: jsonpbPrint,
		ctx:           context.Background(),
	}

	var rootCmd = &cobra.Command{
		Use: "Mesos Operator CLI",
	}
	rootCmd.PersistentFlags().StringSliceVarP(
		&cfg.endpoints, "endpoints", "e", []string{"http://127.0.0.1:5050/api/v1"}, "operator api endpoint",
	)
	rootCmd.AddCommand(agentCommands(cfg))
	rootCmd.AddCommand(masterCommands(cfg))
	rootCmd.Execute()
}

func jsonpbPrint(m proto.Message) {
	marsh := jsonpb.Marshaler{EmitDefaults: true}
	b, err := marsh.MarshalToString(m)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(b)
	}
}
