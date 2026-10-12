package cli

import "github.com/spf13/cobra"

func newStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "gives stat of active windows",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}
