package use

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bf2fc6cc711aee1a0c2a/cli/pkg/api/builders"
	"github.com/bf2fc6cc711aee1a0c2a/cli/pkg/config"
)

func NewUseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [Kafka ID]",
		Short: "Use provided cluster",
		Long:  "Set to work with cluster on current context",
		Args:  cobra.MinimumNArgs(1),
		Run:   runUse,
	}
	return cmd
}

func runUse(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	id := args[0]

	client := builders.BuildClient()

	res, status, err := client.DefaultApi.GetKafkaById(context.Background(), id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving Kafka cluster \"%v\": %v", id, err)
		return
	}

	if status.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Could not use cluster \"%v\": %v", id, err)
		return
	}

	var kafkaConfig config.KafkaConfig = config.KafkaConfig{ClusterID: res.Id, ClusterName: res.Name, ClusterHost: res.BootstrapServerHost}
	cfg.Services.SetKafka(&kafkaConfig)
	if err := config.Save(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Using Kafka cluster \"%v\"", res.Id)
}