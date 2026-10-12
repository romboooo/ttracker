package cli

import (
	"github.com/romboooo/ttracker/tracker"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "track active Hyprland windows",
		Run: func(cmd *cobra.Command, args []string) {
			tracker.Run()
		},
	}
}
