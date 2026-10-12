package cli

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ttracker",
		Short: "tracker for Hyprland active windows",
	}

	cmd.AddCommand(newRunCmd())

	return cmd
}
