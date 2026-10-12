package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/romboooo/ttracker/report"
	"github.com/romboooo/ttracker/storage"
	"github.com/spf13/cobra"
)

func newStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "gives stat of active windows",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := storage.Connect()

			if err != nil {
				return err
			}

			defer store.Close()
			items, err := report.Today(context.Background(), store)
			for _, item := range items {
				fmt.Printf("%s: %s\n", item.Class, item.Duration.Round(time.Second).String())
			}
			return nil
		},
	}
}
