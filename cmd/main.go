package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "br",
		Short: "BartmossRoar Core CLI tool. ",
		Long:  "Hidden your identity, chat free. ",
	}

	rootCmd.AddCommand(joinCmd())
	rootCmd.AddCommand(startCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func joinCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "join",
		Short:   "join a group",
		Long:    "join a bartmoss roar group",
		Example: "br join file.json",
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
}

func startCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start the node",
		Long: `Start the BartmossRoar node.
			The node will:
				1. Discover neighbor nodes in the channel
				2. Sync and rotate ML-KEM public keys
				3. Tick dummy packets at 3-7s intervals
				4. Handle send, gossip and receive`,
		Example: "br start",
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
}
