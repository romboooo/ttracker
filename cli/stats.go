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
		Use:   "stats [today|week|month|year]",
		Short: "gives stat of active windows",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := storage.Connect()

			if err != nil {
				return err
			}

			defer store.Close()

			period := "today"
			if len(args) > 0 {
				period = args[0]
			}

			switch period {
			case "today":
				return printStats(report.Today(context.Background(), store))
			case "week":
				return printStats(report.ThisWeek(context.Background(), store))
			case "month":
				return printStats(report.ThisMonth(context.Background(), store))
			case "year":
				return printStats(report.ThisYear(context.Background(), store))
			default:
				return fmt.Errorf("unknown stats period %q", period)
			}
		},
	}
}

func printStats(items []report.SessionDuration, err error) error {
	if err != nil {
		return err
	}

	for _, item := range items {
		fmt.Printf("%s: %s\n", item.Class, item.Duration.Round(time.Second).String())
	}

	return nil
}
